package main

import (
	"fmt"

	"api.backend.xjco2913/util/config"
	"api.backend.xjco2913/util/zlog"
)

func main() {
	r := NewRouter()

	port := config.Get("server.port")

	zlog.Info(fmt.Sprintf("Starting listening at :%v...", port))
	zlog.Error("this is a error")
	
	r.Run(fmt.Sprintf(":%v", port))
}