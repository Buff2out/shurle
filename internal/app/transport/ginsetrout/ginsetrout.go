package ginsetrout

import (
	"fmt"
	shserv "github.com/Buff2out/shurle/internal/app/services/shurlserv"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

var links = make(map[string]string)

func SetupRouter() *gin.Engine {
	r := gin.Default()

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
		c.String(http.StatusCreated, fmt.Sprintf(`%s:%s%s%s%s%s`, `http`, `/`, `/`, `localhost:8080`, `/`, id))
	})

	return r
}
