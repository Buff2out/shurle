package envs

import (
	"github.com/Buff2out/shurle/internal/app/config/server"
	"github.com/caarlos0/env/v9"
)

type Config struct {
	//Files []string `env:"FILES" envSeparator:":"`
	//Home  string   `env:"HOME"`
	//// required требует, чтобы переменная TASK_DURATION была определена
	//TaskDuration time.Duration `env:"TASK_DURATION,required"`
	Socket string `env:"SERVER_ADDRESS,required"`
	Prefix string `env:"BASE_URL,required"`
}

func GetEnvs() (bool, server.ServerConfig) {
	var cfgparams Config
	//fmt.Println("ошибки нет cfgparams", cfgparams)
	err := env.Parse(&cfgparams)
	cfg := server.ServerConfig{
		S: cfgparams.Socket,
		P: cfgparams.Prefix,
	}
	if err != nil {
		return true, server.ServerConfig{} // true == ДА, МЫ ОШИБЛИСЬ
	}
	//fmt.Println("ошибки нет cfg", cfg)
	return false, cfg // false == НЕТ, ошибка отсутствует, ВСЁ ОК!
}
