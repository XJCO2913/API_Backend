package main

import "api.backend.xjco2913/util/zlog"

func main() {
	r := NewRouter()

	zlog.Info("Starting listening at :8080...")
	
	r.Run(":8080")
}