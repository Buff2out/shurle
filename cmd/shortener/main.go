package main

import (
	pg "github.com/Buff2out/shurle/internal/app/transport/pages"
	"net/http"
)

func main() {
	//fmt.Println("222ffwuef")
	//shello.PrintHello()

	//links := make(map[string]string)
	mux := http.NewServeMux()

	mux.HandleFunc(`/`, pg.HandleShurlPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
