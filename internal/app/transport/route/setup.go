package route

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Buff2out/shurle/internal"
	"github.com/Buff2out/shurle/internal/app/config/server"
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
			panic(err) // TODO заменить панику
		}

		// хм, надо подумать.
		// три строчки ниже - это бизнес логика Save().
		// Можем ли мы сделать некий интерфейс чтобы избежать? дублирования текущего хендлера
		links[id] = string(b)
		eventObj := internal.ShURLFile{UID: strconv.Itoa(len(links)), ShortURL: id, OriginalURL: links[id]}
		filesc.AddNote(sugar, eventObj, filename)

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

		// две строчки ниже - бизнес логика Save()
		eventObj := internal.ShURLFile{UID: strconv.Itoa(len(links)), ShortURL: id, OriginalURL: links[id]}
		filesc.AddNote(sugar, eventObj, filename)
		// формируем ответ
		var respJSON internal.Shlink
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

func Setup(settings *server.Settings, sugar *zap.SugaredLogger) *gin.Engine {
	// версия без контекстов. Потому как у gin framework всё немного
	// усложняется с (c *gin.Context) как пойму как реализовать вместе с ним TODO - переделаю
	DB, errorStartDB := sql.Open("pgx", settings.DatabaseDSN)
	if errorStartDB != nil {

		// Ещё раз нарушаю DRY но исправлю, когда будут
		// одинаковые хендлеры что для случая
		// с открытием бд, что без (через файл и links)
		sugar.Infow("NO CONNECTION DB, got no parameters ", "ERR", errorStartDB, "db", DB)
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
	// комментирую, иначе "DATABASE IS CLOSED",
	// ведь мы выходим из сетап роутера. БД нужно в мейне прописывать"
	// defer DB.Close()
	errExec := repositories.SQLCreateTableURLs(DB)
	if errExec != nil {
		sugar.Infow("INVALID TRY TO CREATE TABLE", "ERR", errExec)
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
	sugar.Infow("SUCCESS TRY TO CREATE TABLE", "ERR", errExec)

	// Отсюда нужно ставить хендлеры которые работают с базой данных
	r := gin.Default()
	r.POST("/", DBPostServeURL(DB, settings.Prefix, sugar))
	r.POST("/:сrutch0/", DBPostServeURL(DB, settings.Prefix, sugar))
	r.POST("/api/shorten", DBPostAPIURL(DB, settings.Prefix, sugar))
	r.GET("/:idvalue", DBGetOriginURL(DB, sugar))

	return r
}
