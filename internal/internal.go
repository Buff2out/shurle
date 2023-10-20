package internal

import (
	"fmt"
	"github.com/Buff2out/shurle/internal/app/config/flags"
	"github.com/caarlos0/env/v9"
	"github.com/fatih/structs"
	"go.uber.org/zap"
)

// НУЖНО ВЫНЕСТИ internal.Settings В shortener!!!

type Settings struct {
	Socket     string `env:"SERVER_ADDRESS,required"`
	Prefix     string `env:"BASE_URL,required"`
	ShURLsJSON string `env:"FILE_STORAGE_PATH,required"`
}

func GetSettings(sugar *zap.SugaredLogger) *Settings {
	var settingsEnvs Settings
	err := env.Parse(&settingsEnvs)
	if err != nil {
		sugar.Infow("GOT ONE OF EMPTY OR NOT EXISTING ENV")
		settingsEnvs = *filterEmptyVals(sugar, &settingsEnvs)
	}
	return &settingsEnvs
}

func filterEmptyVals(sugar *zap.SugaredLogger, settingsEnvs *Settings) *Settings {
	settingsFlags := flags.GetFlags()
	settingsEnvsMap := structs.Map(settingsEnvs)
	settingsFlagsMap := structs.Map(settingsFlags)
	fmt.Println(settingsEnvsMap)
	for key, _ := range settingsEnvsMap {
		if settingsEnvsMap[key] == "" {
			sugar.Infow("Got Key from FLAGS", key, settingsFlagsMap[key])
			settingsEnvsMap[key] = settingsFlagsMap[key]
		}
	}
	return settingsEnvs
}
