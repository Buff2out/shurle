package main

import (
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

	r := ginsetrout.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
