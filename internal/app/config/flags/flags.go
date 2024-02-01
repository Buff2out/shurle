package flags

import (
	"flag"

	"github.com/Buff2out/shurle/internal/app/config/server"
)

var (
	Socket      *string
	Prefix      *string
	ShURLsFile  *string
	DatabaseDSN *string
)

func init() {
	Socket = flag.String("a", `localhost:8080`, "socket = host:port")
	Prefix = flag.String("b", `http://localhost:8080`, "prefix/hashid")
	ShURLsFile = flag.String("f", ``, "Filepath flag")
	DatabaseDSN = flag.String("d", ``, "Database config like host=localhost user=username password=XXXXX dbname=my_db_name sslmode=disable")
}

func GetFlags() *server.Settings {
	flag.Parse()
	return &server.Settings{
		Socket:      *Socket,
		Prefix:      *Prefix,
		ShURLsJSON:  *ShURLsFile,
		DatabaseDSN: *DatabaseDSN,
	}
}
