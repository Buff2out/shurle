package ginsetrout

import (
	"fmt"
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

func MWPostServeUrl(prefix string, sugar *zap.SugaredLogger) func(c *gin.Context) { // mw - не nfs most wanted, а MiddleWare
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

func MWGetOriginUrl(sugar *zap.SugaredLogger) func(c *gin.Context) {
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
	r.GET("/:idvalue", MWGetOriginUrl(sugar))
	r.POST("/", MWPostServeUrl(prefix, sugar))
	r.POST("/:сrutch0/", MWPostServeUrl(prefix, sugar))
	r.POST("/:сrutch0/:сrutch1", MWPostServeUrl(prefix, sugar))
	return r
}
