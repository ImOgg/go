package database

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
	"my-api/config"
)

var RedisClient *redis.Client
var ctx = context.Background()

// 初始化 Redis 連接
func InitRedis() {
	cfg := config.GlobalConfig.Redis
	db, _ := strconv.Atoi(cfg.DB)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       db,
	})

	// 測試連接
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("無法連接到 Redis:", err)
	}

	fmt.Println("Redis 連接成功！")
}

// 獲取 Context
func GetRedisContext() context.Context {
	return ctx
}
