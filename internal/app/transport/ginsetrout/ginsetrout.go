package ginsetrout

import (
	cgzip "compress/gzip"
	"fmt"
	event "github.com/Buff2out/shurle/internal/app/api/shortener"
	"github.com/Buff2out/shurle/internal/app/config/db"
	"github.com/Buff2out/shurle/internal/app/services/filesc"
	"github.com/Buff2out/shurle/internal/app/services/reqsc"
	shserv "github.com/Buff2out/shurle/internal/app/services/shurlsc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
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

		if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
			sugar.Infow(
				"GZIPED request",
			)
			zr, err := cgzip.NewReader(c.Request.Body)
			if err != nil {
				sugar.Infow("Error to create gzipped reader body", "nameErr", err)
			}

			// как в алисе
			c.Request.Body = zr

		}

		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		// здесь я должен предварительно декодировать стрингу, чтобы у нас не хранилось сжатое нами самими.
		links[id] = string(b)
		eventObj := event.ShURLFile{UID: strconv.Itoa(len(links)), ShortURL: id, OriginalURL: links[id]}
		filesc.AddNote(sugar, eventObj, filename)

		// здесь почему-то автотест не ругается, что у меня не предусмотрена возможность отправлять сжатое
		// при Accept-Encoding: gzip. А надо бы.
		c.String(http.StatusCreated, fmt.Sprintf(`%s%s%s`, prefix, `/`, id))
		timeEndingRequest := time.Now()
		sugar.Infow(
			"THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(),
			"StatusCode", strconv.Itoa(http.StatusCreated), // мда, а вот это уже похоже на хардкод, но пусть пока будет так.
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

		if links[id][:4] != "http" {
			links[id] = reqsc.DecodedGzipedOriginURL(links, id)
		}
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

	// Здесь временно (потому что это будет некрасиво, поэтому временно)
	// проинициализируем links из файлов.
	errorStartDB := db.StartDB("pgx", settings.DatabaseDSN)
	links = filesc.FillEvents(sugar, settings.ShURLsJSON, links)
	r := gin.Default()
	r.GET("/ping", MWGetPing(sugar, errorStartDB))
	//r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.GET("/:idvalue", MWGetOriginURL(sugar))
	r.POST("/", MWPostServeURL(settings.Prefix, sugar, settings.ShURLsJSON))
	r.POST("/:сrutch0/", MWPostServeURL(settings.Prefix, sugar, settings.ShURLsJSON))
	r.POST("/:сrutch0/:сrutch1", MWPostServeURL(settings.Prefix, sugar, settings.ShURLsJSON))
	r.POST("/api/shorten", MWPostAPIURL(settings.Prefix, sugar, settings.ShURLsJSON))
	return r
}
