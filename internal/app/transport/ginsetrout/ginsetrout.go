package ginsetrout

import (
	"bytes"
	cgzip "compress/gzip"
	"encoding/json"
	"fmt"
	Event "github.com/Buff2out/shurle/internal/app/api/shortener"
	"github.com/Buff2out/shurle/internal/app/config/files"
	"github.com/Buff2out/shurle/internal/app/services/filesc"
	shserv "github.com/Buff2out/shurle/internal/app/services/shurlsc"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type NetAddress struct {
	Host string
	Port string
}

//type gzreadCloser struct {
//	*cgzip.Reader
//	io.Closer
//}
//
//func (gz gzreadCloser) Close() error {
//	return gz.Closer.Close()
//}

func MWPostServeURL(prefix string, sugar *zap.SugaredLogger, filename string) func(c *gin.Context) { // mw - не nfs most wanted, а MiddleWare
	// наконец-то норм мидлварь, делаем до ретёрна что хотим

	//такс, у меня отличная мысль, реализовать функцию автологгирования, чтобы само время ставилось и кастомное сообщение
	sugar.Infow(
		"POST BLABLA NAH",
	)
	return func(c *gin.Context) {
		timeStartingRequest := time.Now()
		id := shserv.RandStringRunes(5)
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		links[id] = string(b)
		event := Event.ShURLFile{UID: strconv.Itoa(len(links)), ShortURL: id, OriginalURL: links[id]}
		if filename != "" {
			p, er := files.NewProducer(filename, sugar)
			if er != nil {
				sugar.Infow("In MWPostAPIURL func under event var. Invalid path to file.")
			}
			er = p.WriteEvent(&event)
			if er != nil {
				sugar.Infow("In MWPostAPIURL func under WriteEvent. Cant Write to file.")
			}
		}
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
		"POST api url",
	)
	return func(c *gin.Context) {
		timeStartingRequest := time.Now()
		id := shserv.RandStringRunes(5)
		var reqJSON Event.OriginURL
		sugar.Infow(
			"?GZIPED request?", "content-enc", c.Request.Header.Get("Content-Encoding"), "accept-enc", c.Request.Header.Get("Accept-Encoding"),
		)
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

		if err := c.BindJSON(&reqJSON); err != nil {
			sugar.Infow("error in binding json", "nameError", err)
		}
		// Ниже логгируем Json
		//иначе тест не примет
		out, err := json.Marshal(reqJSON)
		if err != nil {
			log.Fatal(err)
		}
		sugar.Infow(
			"json.Unmarshal(b, &reqJSONexmpl)", "reqJSONexmpl = ", out,
		)

		// Записываем хеш в ассоциатор с урлом
		links[id] = reqJSON.URL
		sugar.Infow(
			"reqJSON.URL first 4 symbols", "reqJSON.URL[:4] = ", reqJSON.URL[:4],
		)
		event := Event.ShURLFile{UID: strconv.Itoa(len(links)), ShortURL: id, OriginalURL: links[id]}
		filesc.AddNote(sugar, event, filename)

		// формируем ответ
		var respJSON Event.Shlink
		respJSON.Result = fmt.Sprintf(`%s%s%s`, prefix, `/`, id)
		c.JSON(http.StatusCreated, respJSON)
		timeEndingRequest := time.Now()
		sugar.Infow(
			"THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(),
			"StatusCode", strconv.Itoa(http.StatusCreated),
		)
	}
}

func MWGetOriginURL(sugar *zap.SugaredLogger) func(c *gin.Context) {
	timeStartingServer := time.Now()
	sugar.Infow(
		"StartingServer",
		", Created at", timeStartingServer.String(),
	)
	return func(c *gin.Context) {
		timeStartingRequest := time.Now()

		id := c.Params.ByName("idvalue")

		if links[id][:4] != "http" {
			reader := bytes.NewReader([]byte(links[id]))
			gzreader, e1 := cgzip.NewReader(reader)
			if e1 != nil {
				panic(e1) // Maybe panic here, depends on your error handling.
			}

			output, e2 := io.ReadAll(gzreader)
			if e2 != nil {
				panic(e2)
			}
			links[id] = string(output)
			sugar.Infow(
				"in IF block links[id]", "links[id] = ", links[id],
			)
		}
		sugar.Infow(
			"\"Location\", id", "id = ", id,
			"StatusCode", strconv.Itoa(http.StatusCreated),
		)
		sugar.Infow(
			"\"Location\", links[id]", "links[id] = ", links[id],
			"StatusCode", strconv.Itoa(http.StatusCreated),
		)

		c.Header("Location", links[id])
		c.String(http.StatusTemporaryRedirect, links[id])
		timeEndingRequest := time.Now()
		sugar.Infow(
			"THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(),
			"StatusCode", strconv.Itoa(http.StatusCreated),
		)
	}
}

func fillEvents(sugar *zap.SugaredLogger, file string) {
	//var events = make([]Event.ShURLFile, 0, 5) // мда, теперь это events атавизм
	if file == "" {
		file = filepath.Join(os.TempDir(), "short-url-db.json")
	}
	c, err := files.NewConsumer(file, sugar)
	if err != nil {
		sugar.Infow(
			"in fillEvents failed",
		)
	} else {
		for {
			sugar.Infow(
				"info about path of file", "file", file,
			)
			el, er := c.ReadEvent()
			if er != nil {
				sugar.Infow(
					"END OF FILE", "element", el,
				)
				break
			}
			//events = append(events, *el) // мда, теперь это events атавизм
			links[el.ShortURL] = el.OriginalURL
		}
	}
}

var links = make(map[string]string)

func SetupRouter(settings *Event.Settings, sugar *zap.SugaredLogger) *gin.Engine {

	// Здесь временно (потому что это будет некрасиво, поэтому временно)
	// проинициализируем links из файлов.
	fillEvents(sugar, settings.ShURLsJSON)

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.GET("/:idvalue", MWGetOriginURL(sugar))
	r.POST("/", MWPostServeURL(settings.Prefix, sugar, settings.ShURLsJSON))
	r.POST("/:сrutch0/", MWPostServeURL(settings.Prefix, sugar, settings.ShURLsJSON))
	r.POST("/:сrutch0/:сrutch1", MWPostServeURL(settings.Prefix, sugar, settings.ShURLsJSON))
	r.POST("/api/shorten", MWPostAPIURL(settings.Prefix, sugar, settings.ShURLsJSON))
	return r
}
