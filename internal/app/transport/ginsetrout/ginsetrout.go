package ginsetrout

import (
	"fmt"
	shserv "github.com/Buff2out/shurle/internal/app/services/shurlsc"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type NetAddress struct {
	Host string
	Port string
}

func MWPostServeUrl(prefix string) func(c *gin.Context) { // mw - не nfs most wanted, а MiddleWare
	// наконец-то норм мидлварь, делаем до ретёрна что хотим
	return func(c *gin.Context) {
		id := shserv.EvaluateHashAndReturn()
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		links[id] = string(b)
		c.String(http.StatusCreated, fmt.Sprintf(`%s%s%s`, prefix, `/`, id))
	}
}

func MWGetOriginUrl() func(c *gin.Context) {
	// миддлварь, логгируем что хотим
	return func(c *gin.Context) {
		id := c.Params.ByName("idvalue")

		c.Header("Location", links[id])
		c.String(http.StatusTemporaryRedirect, links[id])
	}
}

var links = make(map[string]string)

func SetupRouter(prefix string) *gin.Engine {
	r := gin.Default()
	r.GET("/:idvalue", MWGetOriginUrl())
	r.POST("/", MWPostServeUrl(prefix))
	r.POST("/:сrutch0/", MWPostServeUrl(prefix))
	r.POST("/:сrutch0/:сrutch1", MWPostServeUrl(prefix))
	return r
}
