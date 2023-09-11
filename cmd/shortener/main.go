package main

import (
	//_ "github.com/Buff2out/shurle/internal/app/config/server"
	"github.com/Buff2out/shurle/internal/app/services/cfgsc"
	"github.com/Buff2out/shurle/internal/app/transport/ginsetrout"
)

func main() {
	////fmt.Println("222ffwuef")
	////shello.PrintHello()
	//
	////links := make(map[string]string)
	//mux := http.NewServeMux()
	//
	//mux.HandleFunc(`/`, pg.HandleShurlPage)
	//
	//err := http.ListenAndServe(`:8080`, mux)
	//if err != nil {
	//	panic(err)
	//}

	//envs.IsEnvsSet()

	serverConfig := cfgsc.GetServerConfig()
	//fmt.Println("ошибки нет main.go:", serverConfig)
	r := ginsetrout.SetupRouter(serverConfig.P)
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(serverConfig.S)
	if err != nil {
		panic(err)
	}

	/*
		package main

		import (
		    "fmt"
		    "net/http"

		    "github.com/gin-gonic/gin"
		)

		func main() {
		    r := gin.Default()
		    r.GET("/foo", func(c *gin.Context) {
		        c.String(http.StatusOK, "bar")
		        fmt.Println(c.Request.Host+c.Request.URL.String())
		    })
		    // Listen and Server in 0.0.0.0:8080
		    r.Run(":8080")
		}
	*/
}
