package reqsc

import (
	"bytes"
	cgzip "compress/gzip"
	"encoding/json"
	Event "github.com/Buff2out/shurle/internal/app/api/shortener"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"log"
	"strings"
)

func GetJSONRequestURL(sugar *zap.SugaredLogger, c *gin.Context) *Event.OriginURL {
	var reqJSON Event.OriginURL
	sugar.Infow(
		"?GZIPED request?", "content-enc", c.Request.Header.Get("Content-Encoding"), "accept-enc", c.Request.Header.Get("Accept-Encoding"),
	)
	if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
		sugar.Infow(
			"GZIPED request",
		)
		zr, err := cgzip.NewReader(c.Request.Body)
		if err != nil {
			sugar.Infow("Error to create gzipped reader body", "nameErr", err)
		}

		// как в алисе
		c.Request.Body = zr

	}

	if err := c.BindJSON(&reqJSON); err != nil {
		sugar.Infow("error in binding json", "nameError", err)
	}
	// Ниже логгируем Json
	//иначе тест не примет
	out, err := json.Marshal(reqJSON)
	if err != nil {
		log.Fatal(err)
	}
	sugar.Infow(
		"json.Unmarshal(b, &reqJSONexmpl)", "reqJSONexmpl = ", out,
	)

	return &reqJSON
}

func DecodedGzipedOriginURL(links map[string]string, id string) string {
	reader := bytes.NewReader([]byte(links[id]))
	gzreader, e1 := cgzip.NewReader(reader)
	if e1 != nil {
		panic(e1)
	}

	output, e2 := io.ReadAll(gzreader)
	if e2 != nil {
		panic(e2)
	}
	return string(output)
}
