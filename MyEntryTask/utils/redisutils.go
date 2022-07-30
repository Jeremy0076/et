package utils

import (
	"fmt"
	"github.com/go-redis/redis"
)

var Rdb *redis.Client

func InitRedis(cfg *RedisConf) error {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	// 默认连接池PoolSize = 10
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       cfg.DB,
		PoolSize: 100,
		//PoolTimeout: 120 * time.Second,
	})

	_, err := Rdb.Ping().Result()
	if err != nil {
		Logs.Error("redis can't ping err :[%v]\n", err)
		return err
	}
	return nil
}

func CloseRedis() {
	if err := Rdb.Close(); err != nil {
		Logs.Error("redis close err: [%v]\n", err)
		return
	}
}
