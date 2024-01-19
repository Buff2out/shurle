package mwhandlers

import (
	"fmt"
	event "github.com/Buff2out/shurle/internal/app/api/shortener"
	"github.com/Buff2out/shurle/internal/app/services/filesc"
	"github.com/Buff2out/shurle/internal/app/services/reqsc"
	shserv "github.com/Buff2out/shurle/internal/app/services/shurlsc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

func PostServeURL(prefix string, sugar *zap.SugaredLogger, filename string) func(c *gin.Context) {
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

func PostAPIURL(prefix string, sugar *zap.SugaredLogger, filename string) func(c *gin.Context) {
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

func GetOriginURL(sugar *zap.SugaredLogger) func(c *gin.Context) {
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