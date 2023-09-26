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

func postServeUrl(c *gin.Context, prefix string) {
	id := shserv.EvaluateHashAndReturn()
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}
	links[id] = string(b)
	c.String(http.StatusCreated, fmt.Sprintf(`%s%s%s`, prefix, `/`, id))
}

func getOriginUrl(c *gin.Context, prefix string) {
	id := c.Params.ByName("idvalue")

	c.Header("Location", links[id])
	c.String(http.StatusTemporaryRedirect, links[id])
}

// ага... Я понял, я делаю роутинг неправильно, вынося за область видимости SetupRouter,
// нужно их по-другому делать
func mwPostServeUrl(prefix string) func(c *gin.Context) { // mw - не nfs most wanted, а MiddleWare
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

func mwGetOriginUrl() func(c *gin.Context) {
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
	/* первый нюанс - когда r.GET, r.POST */
	r.GET("/:idvalue", mwGetOriginUrl())
	r.POST("/", mwPostServeUrl(prefix))

	// Костыли нарушающие DRY для решения .... проблемы автотестера пока так
	r.POST("/:сrutch0/:сrutch1", mwPostServeUrl(prefix))
	r.POST("/:сrutch0/", mwPostServeUrl(prefix))

	return r
}
