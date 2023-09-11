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

var links = make(map[string]string)

func SetupRouter(prefix string) *gin.Engine {
	r := gin.Default()
	/* первый нюанс - когда r.GET, r.POST */
	r.GET("/:idvalue", func(c *gin.Context) {
		id := c.Params.ByName("idvalue")

		c.Header("Location", links[id])
		c.String(http.StatusTemporaryRedirect, links[id])
	})
	r.POST("/", func(c *gin.Context) {
		id := shserv.EvaluateHashAndReturn()
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		links[id] = string(b)
		c.String(http.StatusCreated, fmt.Sprintf(`%s%s%s`, prefix, `/`, id))
	})

	//// Костыли нарушающие DRY для решения .... проблемы автотестера пока так
	//r.POST("/:сrutch0/:сrutch1", func(c *gin.Context) {
	//	id := shserv.EvaluateHashAndReturn()
	//	b, err := io.ReadAll(c.Request.Body)
	//	if err != nil {
	//		panic(err)
	//	}
	//	links[id] = string(b)
	//	c.String(http.StatusCreated, fmt.Sprintf(`%s%s%s`, prefix, `/`, id))
	//})
	//r.POST("/:сrutch0/", func(c *gin.Context) {
	//	id := shserv.EvaluateHashAndReturn()
	//	b, err := io.ReadAll(c.Request.Body)
	//	if err != nil {
	//		panic(err)
	//	}
	//	links[id] = string(b)
	//	c.String(http.StatusCreated, fmt.Sprintf(`%s%s%s`, prefix, `/`, id))
	//})

	return r
}
