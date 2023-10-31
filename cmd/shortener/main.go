package main

import (
	"github.com/Buff2out/shurle/internal"
	lg "github.com/Buff2out/shurle/internal/app/config/logging"
	"github.com/Buff2out/shurle/internal/app/transport/ginsetrout"
	"go.uber.org/zap"
)

func main() {

	/* TODO
	Здесь будут расписаны в дальнейшем задачи для себя
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
	settings := internal.GetSettings(sugar)
	r := ginsetrout.SetupRouter(settings, sugar)
	err := r.Run(settings.Socket)
	if err != nil {
		panic(err)
	}
}
