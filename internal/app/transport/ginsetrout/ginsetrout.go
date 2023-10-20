package ginsetrout

import (
	"bytes"
	cgzip "compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/Buff2out/shurle/internal"
	Event "github.com/Buff2out/shurle/internal/app/api/shortener"
	"github.com/Buff2out/shurle/internal/app/config/files"
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
			p, er := files.NewProducer(filename)
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
				panic(err)
			}

			// как в алисе
			c.Request.Body = zr

		}

		if err := c.BindJSON(&reqJSON); err != nil {
			panic(err)
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
		if filename != "" {
			p, er := files.NewProducer(filename)
			if er != nil {
				sugar.Infow("In MWPostAPIURL func under event var. Invalid path to file.")
			}
			er = p.WriteEvent(&event)
			if er != nil {
				sugar.Infow("In MWPostAPIURL func under WriteEvent. Cant Write to file.")
			}
		}

		// выше логику тоже нужно вынести в функцию...

		// формируем ответ
		var respJSON Event.Shlink
		respJSON.Result = fmt.Sprintf(`%s%s%s`, prefix, `/`, id)
		//out, err := json.MarshalIndent(respJSON, "", "   ")
		//if err != nil {
		//	log.Fatal(err)
		//}
		c.JSON(http.StatusCreated, respJSON)
		timeEndingRequest := time.Now()
		sugar.Infow(
			"THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(),
			"StatusCode", strconv.Itoa(http.StatusCreated),
		)
	}
}

//func MWGetApiURL(sugar *zap.SugaredLogger) func(c *gin.Context) {
//	sugar.Infow(
//		"GET api url",
//	)
//	return func(c *gin.Context) {
//		timeStartingRequest := time.Now()
//		id := c.Params.ByName("idvalue")
//
//		c.Header("Location", links[id])
//		var resp Event.OriginURL
//		resp.URL = links[id]
//		//out, err := json.Marshal(resp)
//		//if err != nil {
//		//	log.Fatal(err)
//		//}
//		c.JSON(http.StatusTemporaryRedirect, resp)
//		timeEndingRequest := time.Now()
//		sugar.Infow(
//			"THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(),
//			"StatusCode", strconv.Itoa(http.StatusCreated), // мда, а вот это уже похоже на хардкод, но пусть пока будет так.
//		)
//	}
//}

// здесь нужно вручную установить нормальный location
func MWGetOriginURL(sugar *zap.SugaredLogger) func(c *gin.Context) {
	// миддлварь, логгируем что хотим
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
			"StatusCode", strconv.Itoa(http.StatusCreated), // мда, а вот это уже похоже на хардкод, но пусть пока будет так.
		)
		sugar.Infow(
			"\"Location\", links[id]", "links[id] = ", links[id],
			"StatusCode", strconv.Itoa(http.StatusCreated), // мда, а вот это уже похоже на хардкод, но пусть пока будет так.
		)

		c.Header("Location", links[id])
		c.String(http.StatusTemporaryRedirect, links[id])
		timeEndingRequest := time.Now()
		sugar.Infow(
			"THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(),
			"StatusCode", strconv.Itoa(http.StatusCreated), // мда, а вот это уже похоже на хардкод, но пусть пока будет так.
		)
	}
}

func fillEvents(sugar *zap.SugaredLogger, file string) {
	var events = make([]Event.ShURLFile, 0, 5) // мда, теперь это events атавизм
	if file == "" {
		file = filepath.Join(os.TempDir(), "short-url-db.json")
	}
	c, err := files.NewConsumer(file)
	if err != nil {
		sugar.Infow(
			"in fillEvents failed",
		)
	}
	// временно заполню линкс отсюда
	for {
		sugar.Infow(
			"in fillEvents failed BEFORE el, er := c.ReadEvent()", "file", file,
		)
		el, er := c.ReadEvent()
		if er != nil {
			sugar.Infow(
				"in fillEvents failed el, er := c.ReadEvent()",
			)
			break
		}
		events = append(events, *el) // мда, теперь это events атавизм
		links[el.ShortURL] = el.OriginalURL
	}
}

var links = make(map[string]string)

func SetupRouter(settings *internal.Settings, sugar *zap.SugaredLogger) *gin.Engine {

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
