package route

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

func DBGetOriginURL(DB *sql.DB, sugar *zap.SugaredLogger) func(c *gin.Context) {
	// Так, стоп нафиг, когда вернусь с обеда, то нужно будет подправить.
	// Ведь я из HTTP GET Запроса получаю hashCode а не из параметров миддлвари.
	return func(c *gin.Context) {
		hashCode := c.Params.ByName("idvalue")
		sugar.Infow("Info about HashID", "id = ", hashCode, "StatusCode", strconv.Itoa(http.StatusCreated))
		// sugar.Infow("Info about OriginURL", "links[id] = ", links[id], "StatusCode", strconv.Itoa(http.StatusCreated))

		var originURLResult string

		row := repositories.SQLGetOriginURL(DB, hashCode)
		// if errRetrieve != nil {
		// 	sugar.Errorw("ginsetrout.DBGetOriginURL() INVALID TRY TO get row from repositories.SQLRetrieveByField(DB, hashCode)", "ERR", errRetrieve)
		// 	return
		// }
		errScan := row.Scan(&originURLResult)
		if errScan != nil {
			sugar.Errorw("ginsetrout.DBGetOriginURL() INVALID SCANNING originURLResult", "ERR", errScan)
			return
		}
		c.Header("Location", originURLResult)
		c.String(http.StatusTemporaryRedirect, originURLResult)
		sugar.Infow("Inside DBGetOriginURL MiddleWare")
	}
}

func DBPostAPIURL(DB *sql.DB, prefix string, sugar *zap.SugaredLogger) func(c *gin.Context) {
	return func(c *gin.Context) {
		timeStartingRequest := time.Now()
		hashCode := shserv.RandStringRunes(5)
		reqJSON := reqsc.GetJSONRequestURL(sugar, c)
		infoURL := shortener.InfoURL{
			OriginalURL: reqJSON.URL,
			ShortURL:    fmt.Sprintf("%s/%s", prefix, hashCode),
			HashCode:    hashCode,
		}
		// TODO В дальнейшем строчку ниже нужно занести в выполнение только при условии,
		// что мы уверены в уникальности записи, в случае если нет,
		// то возвращаем из БД в infoURL существующую запись и статус код другой возвращаем
		errInsert := repositories.SQLInsertURL(DB, &infoURL)
		if errInsert != nil {
			sugar.Infow("errInsert := repositories.SQLInsertURL(DB, &infoURL) ", "ERR", errInsert)
			return
		}
		// eventObj := event.ShURLFile{UID: strconv.Itoa(len(links)), ShortURL: id, OriginalURL: links[id]}
		// filesc.AddNote(sugar, eventObj, filename)

		// формируем ответ
		var respJSON shortener.Shlink
		respJSON.Result = infoURL.ShortURL
		c.JSON(http.StatusCreated, respJSON)
		timeEndingRequest := time.Now()
		sugar.Infow("THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(), "StatusCode", strconv.Itoa(http.StatusCreated))
	}
}

func DBPostServeURL(DB *sql.DB, prefix string, sugar *zap.SugaredLogger) func(c *gin.Context) {
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
