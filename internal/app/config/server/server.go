package server

import "flag"

type ServerConfig struct {
	S string
	P string
}

var (
	Socket *string
	Prefix *string
)

func init() {
	Socket = flag.String("a", `localhost:8080`, "socket = host:port")
	Prefix = flag.String("b", `http://localhost:8080`, "prefix/hashid")
}

func GetServerConfigFromFlags() ServerConfig {

	//serverConfig.Socket = flag.String("a", `localhost:8080`, "socket = host:port")
	//serverConfig.Prefix = flag.String("b", `localhost:8080`, "http://prefix/hashid")
	flag.Parse()
	//// костыль номер бесконечность:
	//if *Socket == "" || *Prefix == "" {
	//	return ServerConfig{
	//		S: "localhost:8080",
	//		P: "http://localhost:8080",
	//	}
	//}
	return ServerConfig{
		S: *Socket,
		P: *Prefix,
	}
}
