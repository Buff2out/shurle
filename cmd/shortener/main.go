package main

import (
	lg "github.com/Buff2out/shurle/internal/app/config/logging"
	"github.com/Buff2out/shurle/internal/app/services/cfgsc"
	"github.com/Buff2out/shurle/internal/app/transport/ginsetrout"
	"go.uber.org/zap"
)

func main() {

	/* TODO
	1 -
	* - метод извлечения (read, консьюм) short-url-db.json >>>--->>> links
		*** определить откуда взять
			а from ENVS
			б from FLAGS
			в from DEFAULTPATH
		READING FILE
	* - метод записывания (write, продьюс) links >>>--->>> short-url-db.json
		*** Определить куда записывать
			A from ENVS
			B from FLAGS
			C DOESN'T WRITE

	2 - заменить паники
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
	err := r.Run(serverConfig.S)
	if err != nil {
		panic(err)
	}
}
