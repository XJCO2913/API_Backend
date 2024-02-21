package redis

import (
	"context"
	"fmt"
	"strconv"

	"api.backend.xjco2913/util/config"
	"api.backend.xjco2913/util/zlog"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	rdb  *redis.Client
)

func init() {
	// Read Redis information from configuration
	host := config.Get("database.redis.host")
	port := config.Get("database.redis.port")
	password := config.Get("database.redis.password")
	dbStr := config.Get("database.redis.db")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		panic(fmt.Sprintf("Invalid Redis DB number: %s", dbStr))
	}

	// Construct Redis connection string
	addr := fmt.Sprintf("%s:%s", host, port)

	// Create Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, 
		DB:       db,
	})

	// test connection
	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic("unable to connect to Redis: " + err.Error())
	}

	// Log Redis connection success
	zlog.Info("Redis initialized successfully", zap.String("address", addr))
}

// RDB returns a singleton of the Redis client
func RDB() *redis.Client {
	return rdb
}
