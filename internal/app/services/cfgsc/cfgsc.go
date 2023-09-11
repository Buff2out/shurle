package cfgsc

import (
	"github.com/Buff2out/shurle/internal/app/config/envs"
	sv "github.com/Buff2out/shurle/internal/app/config/server"
)

func GetServerConfig() sv.ServerConfig {
	isGot, cfg := envs.GetEnvs()
	if isGot {
		return cfg
	}
	return sv.GetServerConfigFromFlags()
}
