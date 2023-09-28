package cfgsc

import (
	"github.com/Buff2out/shurle/internal/app/config/envs"
	sv "github.com/Buff2out/shurle/internal/app/config/server"
	"go.uber.org/zap"
)

func GetServerConfig(sugar *zap.SugaredLogger) sv.ServerConfig {
	isGot, cfg := envs.GetEnvs()
	if isGot { // is gotError
		sugar.Infow(
			"Got ServerConfig From Flags",
		)
		return sv.GetServerConfigFromFlags()
	}
	sugar.Infow(
		"Got ServerConfig From Envs",
	)
	return cfg
}
