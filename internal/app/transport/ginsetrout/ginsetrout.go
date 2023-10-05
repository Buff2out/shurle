package ginsetrout

import (
	cgzip "compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/Buff2out/shurle/internal/app/api/shortener"
	shserv "github.com/Buff2out/shurle/internal/app/services/shurlsc"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

type NetAddress struct {
	Host string
	Port string
}

type gzreadCloser struct {
	*cgzip.Reader
	io.Closer
}

func (gz gzreadCloser) Close() error {
	return gz.Closer.Close()
}

func MWPostServeURL(prefix string, sugar *zap.SugaredLogger) func(c *gin.Context) { // mw - не nfs most wanted, а MiddleWare
	// наконец-то норм мидлварь, делаем до ретёрна что хотим

	//такс, у меня отличная мысль, реализовать функцию автологгирования, чтобы само время ставилось и кастомное сообщение
	sugar.Infow(
		"POST BLABLA NAH",
	)
	return func(c *gin.Context) {
		timeStartingRequest := time.Now()
		id := shserv.GetRandomHash()
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		links[id] = string(b)
		c.String(http.StatusCreated, fmt.Sprintf(`%s%s%s`, prefix, `/`, id))
		timeEndingRequest := time.Now()
		sugar.Infow(
			"THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(),
			"StatusCode", strconv.Itoa(http.StatusCreated), // мда, а вот это уже похоже на хардкод, но пусть пока будет так.
		)
	}
}

func MWPostAPIURL(prefix string, sugar *zap.SugaredLogger) func(c *gin.Context) {
	sugar.Infow(
		"POST api url",
	)
	return func(c *gin.Context) {
		timeStartingRequest := time.Now()
		id := shserv.GetRandomHash()
		//// на деле этот фрагмент кода бессмысленен
		//var reqJSONexmpl shortener.OriginURL
		//b, err := io.ReadAll(c.Request.Body)
		//if err != nil {
		//	panic(err)
		//}
		//if err = json.Unmarshal(b, &reqJSONexmpl); err != nil {
		//	panic(err)
		//}
		//sugar.Infow(
		//	"json.Unmarshal(b, &reqJSONexmpl)", "reqJSONexmpl = ", reqJSONexmpl,
		//	"b = ", string(b),
		//	"StatusCode", strconv.Itoa(http.StatusCreated), // мда, а вот это уже похоже на хардкод, но пусть пока будет так.
		//)
		//// конец фрагмента

		// А вот этот НИЖЕ - ключевой фрагмент, в котором используется JSON.
		var reqJSON shortener.OriginURL
		sugar.Infow(
			"?GZIPED request?", "content-enc", c.GetHeader("Content-Encoding"),
		)
		if c.GetHeader("Content-Encoding") == "gzip" {
			sugar.Infow(
				"GZIPED request",
			)
			zr, err := cgzip.NewReader(c.Request.Body)
			if err != nil {
				panic(err)
			}
			// потом вынесу это в отдельный метод/функ
			//remember DRY
			b, err := io.ReadAll(zr)
			if err != nil {
				panic(err)
			}
			if err = json.Unmarshal(b, &reqJSON); err != nil {
				panic(err)
			}
			c.Request.Body = gzreadCloser{zr, c.Request.Body}

		} else {
			// потом вынесу это в отдельный метод/функ
			//remember DRY
			sugar.Infow(
				"NOT GZIPED request",
			)
			b, err := io.ReadAll(c.Request.Body)
			if err != nil {
				panic(err)
			}
			if err = json.Unmarshal(b, &reqJSON); err != nil {
				panic(err)
			}
		}

		//if err := c.BindJSON(&reqJSON); err != nil {
		//	panic(err)
		//}
		//// Ниже логгируем Json иначе тест не примет
		//out, err := json.Marshal(reqJSON)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//sugar.Infow(
		//	"json.Unmarshal(b, &reqJSONexmpl)", "reqJSONexmpl = ", out,
		//)

		// Записываем хеш в ассоциатор с урлом
		links[id] = reqJSON.URL

		// формируем ответ
		var respJSON shortener.Shlink
		respJSON.Result = fmt.Sprintf(`%s%s%s`, prefix, `/`, id)
		//out, err := json.MarshalIndent(respJSON, "", "   ")
		//if err != nil {
		//	log.Fatal(err)
		//}
		c.JSON(http.StatusCreated, respJSON)
		timeEndingRequest := time.Now()
		sugar.Infow(
			"THIS IS A REQUEST RESPONSE LOG", "Request duration", timeStartingRequest.Sub(timeEndingRequest).String(),
			"StatusCode", strconv.Itoa(http.StatusCreated), // мда, а вот это уже похоже на хардкод, но пусть пока будет так.
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
//		var resp shortener.OriginURL
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

var links = make(map[string]string)

func SetupRouter(prefix string, sugar *zap.SugaredLogger) *gin.Engine {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.GET("/:idvalue", MWGetOriginURL(sugar))
	r.POST("/", MWPostServeURL(prefix, sugar))
	r.POST("/:сrutch0/", MWPostServeURL(prefix, sugar))
	r.POST("/:сrutch0/:сrutch1", MWPostServeURL(prefix, sugar))
	r.POST("/api/shorten", MWPostAPIURL(prefix, sugar))
	return r
}
