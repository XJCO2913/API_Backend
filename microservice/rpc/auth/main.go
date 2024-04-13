package main

import (
	auth "rpc.xjco2913/kitex_gen/rpc/xjco2913/auth/loginservice"
	"log"
)

func main() {
	svr := auth.NewServer(new(LoginServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
