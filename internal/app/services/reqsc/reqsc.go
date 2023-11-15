package reqsc

import (
	cgzip "compress/gzip"
	"encoding/json"
	Event "github.com/Buff2out/shurle/internal/app/api/shortener"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"log"
	"strings"
)

func DecompressedGZReader(sugar *zap.SugaredLogger, c *gin.Context) (io.ReadCloser, error, string) {
	if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
		sugar.Infow(
			"GZIPED request",
		)
		zr, err := cgzip.NewReader(c.Request.Body)
		if err != nil {
			sugar.Infow("Error to create gzipped reader body", "nameErr", err)
			return nil, err, ""
		}

		// как в алисе
		return zr, nil, "gzip"
	}
	// Опционально можно масштабировать данную функцию, если вдруг есть другие Content-Encoding
	return c.Request.Body, nil, "default"
}

func GetJSONRequestURL(sugar *zap.SugaredLogger, c *gin.Context) *Event.OriginURL {
	var reqJSON Event.OriginURL
	sugar.Infow(
		"?GZIPED request?", "content-enc", c.Request.Header.Get("Content-Encoding"), "accept-enc", c.Request.Header.Get("Accept-Encoding"),
	)
	var err error
	var enc string
	c.Request.Body, err, enc = DecompressedGZReader(sugar, c)
	if err != nil {
		sugar.Infow("Error to create gzipped reader body in GetJSONRequestURL", "nameErr", err)
	}

	if err = c.BindJSON(&reqJSON); err != nil {
		sugar.Infow("error in binding json", "nameError", err)
	}
	//reqJSON.URL = DecodedStringWithEncodingType(sugar, enc, reqJSON.URL)
	// Ниже логгируем Json
	//иначе тест не примет
	out, err := json.Marshal(reqJSON)
	if err != nil {
		log.Fatal(err)
	}
	sugar.Infow(
		"json.Unmarshal(b, &reqJSONexmpl)", "reqJSONexmpl = ", out, "encoding", enc,
	)

	return &reqJSON
}

//func DecodedGzipedOriginURL(links map[string]string, id string) string {
//	reader := bytes.NewReader([]byte(links[id]))
//	gzreader, e1 := cgzip.NewReader(reader)
//	if e1 != nil {
//		panic(e1)
//	}
//
//	output, e2 := io.ReadAll(gzreader)
//	if e2 != nil {
//		panic(e2)
//	}
//	return string(output)
//}

//func DecodedStringWithEncodingType(sugar *zap.SugaredLogger, enc string, str string) string {
//
//	switch enc {
//	case "default":
//		return str
//	case "gzip":
//		reader := bytes.NewReader([]byte(str))
//		gzreader, e1 := cgzip.NewReader(reader)
//		if e1 != nil {
//			// пока что лень паники переделывать под return Ошибок, так пусть пока
//			// порабоает, лучше сфокусироваться на функционале
//			panic(e1)
//		}
//		output, e2 := io.ReadAll(gzreader)
//		if e2 != nil {
//			panic(e2)
//		}
//		return string(output)
//	}
//	sugar.Infow("Error In DecodedStringWithEncodingType! THERE IS NO CASE TYPE FOUNDED")
//	return ""
//}
