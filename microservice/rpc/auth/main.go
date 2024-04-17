package main

import (
	auth "api.backend.xjco2913/microservice/kitex_gen/rpc/xjco2913/auth/authservice"
	"log"
)

func main() {
	svr := auth.NewServer(new(AuthServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
