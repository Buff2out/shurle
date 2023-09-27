package main

import (
	"github.com/Buff2out/shurle/internal/app/services/cfgsc"
	"github.com/Buff2out/shurle/internal/app/transport/ginsetrout"
)

func main() {

	/* TODO
	1 - Навести порядочек в "контроллерах"
	2 - осознать предложения правок Владислава (своего ментора)
	3 - причесать, убрать ненужные комментарии с устаревшим закоменченным кодом
	4 - Параллельно с этим делать инкремент
	*/

	serverConfig := cfgsc.GetServerConfig()
	//fmt.Println("ошибки нет main.go:", serverConfig)
	r := ginsetrout.SetupRouter(serverConfig.P)
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(serverConfig.S)
	if err != nil {
		panic(err)
	}
}
