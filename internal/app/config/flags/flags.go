package flags

import (
	"flag"
	"github.com/Buff2out/shurle/internal/app/api/shortener"
)

// НУЖНО ВЫНЕСТИ internal.Settings В shortener!!!

var (
	Socket     *string
	Prefix     *string
	ShURLsFile *string
)

func init() {
	Socket = flag.String("a", `localhost:8080`, "socket = host:port")
	Prefix = flag.String("b", `http://localhost:8080`, "prefix/hashid")
	ShURLsFile = flag.String("f", ``, "Filepath flag")
}

func GetFlags() *shortener.Settings {
	flag.Parse()
	return &shortener.Settings{
		Socket:     *Socket,
		Prefix:     *Prefix,
		ShURLsJSON: *ShURLsFile,
	}
}
