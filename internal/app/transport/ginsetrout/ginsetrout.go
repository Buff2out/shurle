package ginsetrout

import (
	"encoding/json"
	"fmt"
	"github.com/Buff2out/shurle/internal/app/api/shortener"
	shserv "github.com/Buff2out/shurle/internal/app/services/shurlsc"
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

func MWPostApiURL(prefix string, sugar *zap.SugaredLogger) func(c *gin.Context) {
	sugar.Infow(
		"POST api url",
	)
	return func(c *gin.Context) {
		timeStartingRequest := time.Now()
		id := shserv.GetRandomHash()
		b, err := io.ReadAll(c.Request.Body)

		var reqJSON shortener.OriginURL
		var respJSON shortener.Shlink

		if err != nil {
			panic(err)
		}
		if err = json.Unmarshal(b, &reqJSON); err != nil {
			panic(err)
		}
		// Записываем хеш в ассоциатор с урлом
		links[id] = reqJSON.URL

		// формируем ответ
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
	r.GET("/:idvalue", MWGetOriginURL(sugar))
	r.POST("/", MWPostServeURL(prefix, sugar))
	r.POST("/:сrutch0/", MWPostServeURL(prefix, sugar))
	r.POST("/:сrutch0/:сrutch1", MWPostServeURL(prefix, sugar))
	r.POST("/api/shorten", MWPostApiURL(prefix, sugar))
	return r
}
