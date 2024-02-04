package reqsc

import (
	cgzip "compress/gzip"
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/Buff2out/shurle/internal"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func DecompressedGZReader(sugar *zap.SugaredLogger, c *gin.Context) (io.ReadCloser, string, error) {
	if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
		sugar.Infow(
			"GZIPED request",
		)
		zr, err := cgzip.NewReader(c.Request.Body)
		if err != nil {
			sugar.Infow("Error to create gzipped reader body", "nameErr", err)
			return nil, "", err
		}

		// как в алисе
		return zr, "gzip", nil
	}
	// Опционально можно масштабировать данную функцию, если вдруг есть другие Content-Encoding
	return c.Request.Body, "default", nil
}

func GetJSONRequestURL(sugar *zap.SugaredLogger, c *gin.Context) *internal.OriginURL {
	var reqJSON internal.OriginURL
	sugar.Infow(
		"?GZIPED request?", "content-enc", c.Request.Header.Get("Content-Encoding"), "accept-enc", c.Request.Header.Get("Accept-Encoding"),
	)
	var err error
	var enc string
	c.Request.Body, enc, err = DecompressedGZReader(sugar, c)
	if err != nil {
		sugar.Infow("Error to create gzipped reader body in GetJSONRequestURL", "nameErr", err)
	}

	if err = c.BindJSON(&reqJSON); err != nil {
		sugar.Infow("error in binding json", "nameError", err)
	}
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
