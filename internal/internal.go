package internal

import (
	"fmt"

	"github.com/Buff2out/shurle/internal/app/config/flags"
	"github.com/Buff2out/shurle/internal/app/config/server"
	"github.com/caarlos0/env/v9"
	"github.com/fatih/structs"
	"go.uber.org/zap"
)

type OriginURL struct {
	URL string `json:"url"`
}

type Shlink struct {
	Result string `json:"result,omitempty"`
}

type ShURLFile struct {
	UID         string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type InfoURL struct {
	ShortURL    string
	OriginalURL string
	HashCode    string
}

func GetSettings(sugar *zap.SugaredLogger) *server.Settings {
	var settingsEnvs server.Settings
	err := env.Parse(&settingsEnvs)
	if err != nil {
		return filterEmptyVals(sugar, &settingsEnvs)
	}
	return &settingsEnvs
}

func filterEmptyVals(sugar *zap.SugaredLogger, settingsEnvs *server.Settings) *server.Settings {
	settingsEnvsMap := structs.Map(settingsEnvs)
	settingsFlagsMap := structs.Map(flags.GetFlags())
	for key := range settingsEnvsMap {
		if settingsEnvsMap[key] == "" {
			sugar.Infow("Got Key from FLAGS", key, settingsFlagsMap[key], "previosVal", settingsEnvsMap[key])
			settingsEnvsMap[key] = settingsFlagsMap[key]
		}
	}
	sugar.Infow("Settings", "DatabaseDSN", settingsEnvsMap["DatabaseDSN"])
	return &server.Settings{
		Socket:      fmt.Sprintf("%v", settingsEnvsMap["Socket"]),
		Prefix:      fmt.Sprintf("%v", settingsEnvsMap["Prefix"]),
		ShURLsJSON:  fmt.Sprintf("%v", settingsEnvsMap["ShURLsJSON"]),
		DatabaseDSN: fmt.Sprintf("%v", settingsEnvsMap["DatabaseDSN"]),
	}
}
