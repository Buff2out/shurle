package main

import (
	lg "github.com/Buff2out/shurle/internal/app/config/logging"
	"github.com/Buff2out/shurle/internal/app/services/cfgsc"
	"github.com/Buff2out/shurle/internal/app/transport/ginsetrout"
	"go.uber.org/zap"
)

func main() {
	// comment for testcommit
	/* TODO
	1 - Навести порядочек в "контроллерах"
	2 - осознать предложения правок Владислава (своего ментора)
	3 - причесать, убрать ненужные комментарии с устаревшим закоменченным кодом
	4 - Параллельно с этим делать инкремент
	*/
	sugar, logger := lg.GetSugaredLogger()
	// это нужно добавить, если логер буферизован
	// в данном случае не буферизован, но привычка хорошая
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic("cannot close zap's sugared logger")
		}
	}(logger)
	serverConfig := cfgsc.GetServerConfig(sugar)
	//fmt.Println("ошибки нет main.go:", serverConfig)
	r := ginsetrout.SetupRouter(serverConfig.P, sugar)
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(serverConfig.S)
	if err != nil {
		panic(err)
	}
}
