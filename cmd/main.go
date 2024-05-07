package main

import (
	"context"
	"fmt"
	"os"

	"api.backend.xjco2913/util/config"
	"api.backend.xjco2913/util/zlog"
)

var (
	port string
)

func main() {
	ctx := context.Background()

	ClearJunkAvatarFile(ctx)
	ClearJunkMomentMediaFile(ctx)
	ClearJunkActivityImageFile(ctx)

	r := NewRouter()

	go localHub.Run()

	port = "8080"
	if env, ok := os.LookupEnv("DEPLOY_ENV"); ok {
		if env == "test" {
			port = config.Get("server.test.port")
		} else {
			port = config.Get("server.live.port")
		}
	}

	// Async flush logs into mysql
	go FlushLogs(ctx)

	// Clear junk avatar object file

	zlog.Info(fmt.Sprintf("Starting listening at :%v...", port))
	r.Run(fmt.Sprintf(":%v", port))
}
