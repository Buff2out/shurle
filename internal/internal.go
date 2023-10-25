package internal

import (
	"fmt"
	"github.com/Buff2out/shurle/internal/app/api/shortener"
	"github.com/Buff2out/shurle/internal/app/config/flags"
	"github.com/caarlos0/env/v9"
	"github.com/fatih/structs"
	"go.uber.org/zap"
)

func GetSettings(sugar *zap.SugaredLogger) *shortener.Settings {
	var settingsEnvs shortener.Settings
	err := env.Parse(&settingsEnvs)
	if err != nil {
		return filterEmptyVals(sugar, &settingsEnvs)
	}
	return &settingsEnvs
}

func filterEmptyVals(sugar *zap.SugaredLogger, settingsEnvs *shortener.Settings) *shortener.Settings {
	settingsEnvsMap := structs.Map(settingsEnvs)
	settingsFlagsMap := structs.Map(flags.GetFlags())
	for key := range settingsEnvsMap {
		if settingsEnvsMap[key] == "" {
			sugar.Infow("Got Key from FLAGS", key, settingsFlagsMap[key])
			settingsEnvsMap[key] = settingsFlagsMap[key]
		}
	}
	return &shortener.Settings{
		Socket:     fmt.Sprintf("%v", settingsEnvsMap["Socket"]),
		Prefix:     fmt.Sprintf("%v", settingsEnvsMap["Prefix"]),
		ShURLsJSON: fmt.Sprintf("%v", settingsEnvsMap["ShURLsJSON"]),
	}
}
