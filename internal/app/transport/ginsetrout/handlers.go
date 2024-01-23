package ginsetrout

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Buff2out/shurle/internal/app/api/shortener"
	"github.com/Buff2out/shurle/internal/app/repositories"
	"github.com/Buff2out/shurle/internal/app/services/reqsc"
	shserv "github.com/Buff2out/shurle/internal/app/services/shurlsc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*
сюда добавим хендлеры пакаги для задания джинроутинга
Пока что такие

DBPostDBServeURL
DBPostServeURL
DBPostAPIURL
*/

func DBGetOriginURL(sugar *zap.SugaredLogger) func(c *gin.Context) {

	return func(c *gin.Context) {
		sugar.Infow("Inside DBGetOriginURL MiddleWare")
	}
}

func DBPostServeURL(DB *sql.DB, prefix string, sugar *zap.SugaredLogger) func(c *gin.Context) {
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

		// три строчки ниже - это бизнес логика Save().
		infoURL := shortener.InfoURL{
			OriginalURL: string(b),
			ShortURL:    fmt.Sprintf("%s/%s", prefix, id),
			HashCode:    id,
		}
		// TODO В дальнейшем строчку ниже нужно занести в выполнение только при условии,
		// что мы уверены в уникальности записи, в случае если нет,
		// то возвращаем из БД в infoURL существующую запись и статус код другой возвращаем
		errInsert := repositories.SQLInsertURL(DB, &infoURL)
		if errInsert != nil {
			sugar.Infow("errInsert := repositories.SQLInsertURL(DB, &infoURL) ", "ERR", errInsert)
		}
		// eventObj := event.ShURLFile{UID: strconv.Itoa(len(links)), ShortURL: id, OriginalURL: links[id]}
		// filesc.AddNote(sugar, eventObj, filename)

		c.String(http.StatusCreated, infoURL.ShortURL)
		timeEndingRequest := time.Now()
		sugar.Infow(
			"THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(),
			"StatusCode", strconv.Itoa(http.StatusCreated), "encoding", enc,
		)
	}
}

func DBPostAPIURL(prefix string, sugar *zap.SugaredLogger, filename string) func(c *gin.Context) {
	return nil
}
