package internal

import (
	"fmt"

	"github.com/Buff2out/shurle/internal/app/config/flags"
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

type Settings struct {
	Socket      string `env:"SERVER_ADDRESS,required"`
	Prefix      string `env:"BASE_URL,required"`
	ShURLsJSON  string `env:"FILE_STORAGE_PATH,required"`
	DatabaseDSN string `env:"DATABASE_DSN"`
}

func GetSettings(sugar *zap.SugaredLogger) *Settings {
	var settingsEnvs Settings
	err := env.Parse(&settingsEnvs)
	if err != nil {
		return filterEmptyVals(sugar, &settingsEnvs)
	}
	return &settingsEnvs
}

func filterEmptyVals(sugar *zap.SugaredLogger, settingsEnvs *Settings) *Settings {
	settingsEnvsMap := structs.Map(settingsEnvs)
	settingsFlagsMap := structs.Map(flags.GetFlags())
	for key := range settingsEnvsMap {
		if settingsEnvsMap[key] == "" {
			sugar.Infow("Got Key from FLAGS", key, settingsFlagsMap[key], "previosVal", settingsEnvsMap[key])
			settingsEnvsMap[key] = settingsFlagsMap[key]
		}
	}
	sugar.Infow("Settings", "DatabaseDSN", settingsEnvsMap["DatabaseDSN"])
	return &Settings{
		Socket:      fmt.Sprintf("%v", settingsEnvsMap["Socket"]),
		Prefix:      fmt.Sprintf("%v", settingsEnvsMap["Prefix"]),
		ShURLsJSON:  fmt.Sprintf("%v", settingsEnvsMap["ShURLsJSON"]),
		DatabaseDSN: fmt.Sprintf("%v", settingsEnvsMap["DatabaseDSN"]),
	}
}
