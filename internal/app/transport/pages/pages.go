package pages

// this file is redundant

import (
	"fmt"
	shserv "github.com/Buff2out/shurle/internal/app/services/shurlsc"
	"io"
	"net/http"
	"strings"
)

var links = make(map[string]string)

func HandleShurlPage(res http.ResponseWriter, req *http.Request) {
	//res.Header().Set("Access-Control-Allow-Origin", "*")
	if http.MethodPost == req.Method { // Добавить ещё условие проверки длинности

		res.Header().Set("content-type", "text/plain") // wow, http.ContentTypeText doesn't work

		res.WriteHeader(http.StatusCreated) // 201
		hash := shserv.RandStringRunes(5)

		b, err := io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		links[hash] = string(b)
		_, err = res.Write([]byte(fmt.Sprintf(`%s:%s%s%s%s%s`, `http`, `/`, `/`, `localhost:8080`, `/`, hash)))
		if err != nil {
			return
		}
	} else if http.MethodGet == req.Method {
		//fmt.Println("gsdf")
		res.Header().Set("content-type", "text/plain") // wow, http.ContentTypeText doesn't work
		hash := strings.Split(req.URL.Path, "/")[1]
		res.Header().Set("Location", (fmt.Sprintf(links[hash])))
		res.WriteHeader(http.StatusTemporaryRedirect) // 307
		//_, err := res.Write([]byte(fmt.Sprintf(`http://%s`, links[hash])))
		//if err != nil {
		//	return
		//}
	} else {
		res.WriteHeader(http.StatusBadRequest) // 201
	}
}

//func GetLinkByHashAndRedirect(res http.ResponseWriter, req *http.Request) {
//	res.Header().Set("Access-Control-Allow-Origin", "*")
//	res.Header().Set("content-type", "text/plain")
//	res.WriteHeader(http.StatusTemporaryRedirect) // 307
//
//}
