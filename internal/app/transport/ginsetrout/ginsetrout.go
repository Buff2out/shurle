package ginsetrout

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	event "github.com/Buff2out/shurle/internal/app/api/shortener"
	"github.com/Buff2out/shurle/internal/app/repositories"
	"github.com/Buff2out/shurle/internal/app/services/filesc"
	"github.com/Buff2out/shurle/internal/app/services/reqsc"
	shserv "github.com/Buff2out/shurle/internal/app/services/shurlsc"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type NetAddress struct {
	Host string
	Port string
}

func MWPostServeURL(prefix string, sugar *zap.SugaredLogger, filename string) func(c *gin.Context) {
	sugar.Infow(
		"Info message inside MWPostServeURL",
	)
	return func(c *gin.Context) {
		timeStartingRequest := time.Now()
		id := shserv.RandStringRunes(5)
		var err error
		var enc string
		c.Request.Body, enc, err = reqsc.DecompressedGZReader(sugar, c)
		if err != nil {
			sugar.Infow("Error to create gzipped reader body", "nameErr", err)
		}

		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		//links[id] = reqsc.DecodedStringWithEncodingType(sugar, enc, string(b))
		links[id] = string(b)
		eventObj := event.ShURLFile{UID: strconv.Itoa(len(links)), ShortURL: id, OriginalURL: links[id]}
		filesc.AddNote(sugar, eventObj, filename)

		// здесь почему-то автотест не ругается, что у меня не предусмотрена возможность отправлять сжатое
		// при Accept-Encoding: gzip. А надо бы.
		c.String(http.StatusCreated, fmt.Sprintf(`%s%s%s`, prefix, `/`, id))
		timeEndingRequest := time.Now()
		sugar.Infow(
			"THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(),
			"StatusCode", strconv.Itoa(http.StatusCreated), "encoding", enc,
		)
	}
}

func MWPostAPIURL(prefix string, sugar *zap.SugaredLogger, filename string) func(c *gin.Context) {
	sugar.Infow(
		"Info message inside MWPostAPIURL",
	)
	return func(c *gin.Context) {
		timeStartingRequest := time.Now()
		id := shserv.RandStringRunes(5)
		reqJSON := reqsc.GetJSONRequestURL(sugar, c)
		// Записываем хеш в ассоциатор с урлом
		links[id] = reqJSON.URL
		sugar.Infow("reqJSON.URL first 4 symbols", "reqJSON.URL[:4] = ", reqJSON.URL[:4])
		eventObj := event.ShURLFile{UID: strconv.Itoa(len(links)), ShortURL: id, OriginalURL: links[id]}
		filesc.AddNote(sugar, eventObj, filename)
		// формируем ответ
		var respJSON event.Shlink
		respJSON.Result = fmt.Sprintf(`%s%s%s`, prefix, `/`, id)
		c.JSON(http.StatusCreated, respJSON)
		timeEndingRequest := time.Now()
		sugar.Infow("THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(), "StatusCode", strconv.Itoa(http.StatusCreated))
	}
}

func MWGetOriginURL(sugar *zap.SugaredLogger) func(c *gin.Context) {
	timeStartingServer := time.Now()
	sugar.Infow("StartingServer", "Created at", timeStartingServer.String())
	return func(c *gin.Context) {
		timeStartingRequest := time.Now()

		id := c.Params.ByName("idvalue")
		sugar.Infow("Info about HashID", "id = ", id, "StatusCode", strconv.Itoa(http.StatusCreated))
		sugar.Infow("Info about OriginURL", "links[id] = ", links[id], "StatusCode", strconv.Itoa(http.StatusCreated))

		c.Header("Location", links[id])
		c.String(http.StatusTemporaryRedirect, links[id])
		timeEndingRequest := time.Now()
		sugar.Infow(
			"THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(),
			"StatusCode", strconv.Itoa(http.StatusCreated),
		)
	}
}

func MWGetPing(sugar *zap.SugaredLogger, errorStartDB error) func(c *gin.Context) {
	timeStartingServer := time.Now()
	sugar.Infow("StartingServer", "Created at", timeStartingServer.String())
	return func(c *gin.Context) {
		timeStartingRequest := time.Now()
		if errorStartDB != nil {
			sugar.Errorw("Error starting db", "texterr", errorStartDB)
			c.String(http.StatusInternalServerError, "")
			return
		}
		c.String(http.StatusOK, "")
		timeEndingRequest := time.Now()
		sugar.Infow(
			"THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(),
			"StatusCode", strconv.Itoa(http.StatusCreated),
		)
	}
}

var links = make(map[string]string)

func SetupRouter(settings *event.Settings, sugar *zap.SugaredLogger) *gin.Engine {
	// версия без контекстов. Потому как у gin framework всё немного
	// усложняется с (c *gin.Context) как пойму как реализовать вместе с ним TODO - переделаю
	urlsDB, errorStartDB := sql.Open("pgx", settings.DatabaseDSN)
	if errorStartDB != nil {
		sugar.Infow("NO CONNECTION DB, got no parameters ", "ERR", errorStartDB, "db", urlsDB)
		links = filesc.FillEvents(sugar, settings.ShURLsJSON, links)
		r := gin.Default()
		// r.GET("/ping", MWGetPing(sugar, errorStartDB))
		r.GET("/:idvalue", MWGetOriginURL(sugar))
		r.POST("/", MWPostServeURL(settings.Prefix, sugar, settings.ShURLsJSON))
		r.POST("/:сrutch0/", MWPostServeURL(settings.Prefix, sugar, settings.ShURLsJSON))
		r.POST("/:сrutch0/:сrutch1", MWPostServeURL(settings.Prefix, sugar, settings.ShURLsJSON))
		r.POST("/api/shorten", MWPostAPIURL(settings.Prefix, sugar, settings.ShURLsJSON))
		return r
	}
	defer urlsDB.Close()
	errExec := repositories.SQLCreateTableURLs(urlsDB)
	if errExec != nil {
		sugar.Infow("INVALID TRY TO EXEC", "ERR", errExec)
		return gin.Default()
	}
	sugar.Infow("SUCCESS TRY TO EXEC", "ERR", errExec)
	errInsert := repositories.SQLInsert(urlsDB)
	if errInsert != nil {
		sugar.Infow("INVALID TRY TO INSERT", "ERR", errInsert)
		return gin.Default()
	}
	sugar.Infow("SUCCESS TRY TO INSERT", "ERR", errInsert)
	// наконец то придумал как реализовать обратную совместимость и
	// и впихнуть логику с бд как только стало ясно, что
	// в переменных окружения получилось получить DSN
	// мы просто пропишем отдельные хендлеры для них
	// а тут сделаем проверку на то, есть ли тут databaseDSN.

	// links = filesc.FillEvents(sugar, settings.ShURLsJSON, links)
	r := gin.Default()
	// r.GET("/ping", MWGetPing(sugar, errorStartDB))
	// r.GET("/:idvalue", DBGetDBOriginURL(sugar))
	// r.POST("/", DBPostDBServeURL(settings.Prefix, sugar, settings.ShURLsJSON))
	// r.POST("/:сrutch0/", DBPostDBServeURL(settings.Prefix, sugar, settings.ShURLsJSON))
	// r.POST("/:сrutch0/:сrutch1", DBPostServeURL(settings.Prefix, sugar, settings.ShURLsJSON))
	// r.POST("/api/shorten", DBPostAPIURL(settings.Prefix, sugar, settings.ShURLsJSON))
	return r
}
