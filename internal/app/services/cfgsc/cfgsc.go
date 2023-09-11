package cfgsc

import (
	"fmt"
	"github.com/Buff2out/shurle/internal/app/config/envs"
	sv "github.com/Buff2out/shurle/internal/app/config/server"
)

func GetServerConfig() sv.ServerConfig {
	fmt.Println("GOT IT 0")
	isGot, cfg := envs.GetEnvs()
	if isGot {
		fmt.Println("GOT IT 1")
		return cfg
	}
	fmt.Println("GOT IT 2")
	return sv.GetServerConfigFromFlags()
}
