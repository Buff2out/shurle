package cfgsc

import (
	"github.com/Buff2out/shurle/internal/app/config/envs"
	sv "github.com/Buff2out/shurle/internal/app/config/server"
	"go.uber.org/zap"
)

// TODO
// ПЕРЕНЕСТИ ЭТОТ МЕТОД ИЗ СЕРВИСОВ В CONFIG.GO
// ТАМ СДЕЛАТЬ НЕКИЙ SETUP CONFIG METHOD
func GetServerConfig(sugar *zap.SugaredLogger) sv.ServerConfig {
	isGot, cfg := envs.GetEnvs()
	if isGot { // is gotError (upd: 2023.10.06 - читаю этот бред сейчас и понимаю,
		// насколько тяжело это читать, переделаем скоро.)
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
