package main

import (
	"fmt"
	"os"

	"api.backend.xjco2913/util/config"
	"api.backend.xjco2913/util/zlog"
)

var (
	port string
)

func main() {
	r := NewRouter()

	if env, ok := os.LookupEnv("DEPLOY_ENV"); ok {
		if env == "test" {
			port = config.Get("server.test.port")
		} else {
			port = config.Get("server.live.port")
		}
	}

	zlog.Info(fmt.Sprintf("Starting listening at :%v...", port))
	zlog.Error("this is a error")
	
	r.Run(fmt.Sprintf(":%v", port))
}