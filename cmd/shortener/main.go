package main

import (
	"github.com/Buff2out/shurle/internal"
	lg "github.com/Buff2out/shurle/internal/app/config/logging"
	"github.com/Buff2out/shurle/internal/app/transport/ginsetrout"
	"go.uber.org/zap"
)

func main() {
	// comment for testcommit
	/* TODO
	1 -
		проверить если файла не существует, создать новый,
		проэкспериментировать в локальном плейграунде
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
	//serverConfig := cfgsc.GetServerConfig(sugar)
	settings := internal.GetSettings(sugar)
	//fmt.Println("ошибки нет main.go:", serverConfig)
	r := ginsetrout.SetupRouter(settings, sugar)
	err := r.Run(settings.Socket)
	if err != nil {
		panic(err)
	}
}
