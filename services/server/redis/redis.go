package redis

import (
	"streamserver/env"
	log "streamserver/log"

	"github.com/go-redis/redis"
	"github.com/nitishm/go-rejson"
)

var Client *rejson.Handler
var rdb *redis.Client

func Connect() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     env.Vars.Redis.Address,
		Password: env.Vars.Redis.Password,
		DB:       env.Vars.Redis.DB,
	})

	rh := rejson.NewReJSONHandler()

	rh.SetGoRedisClient(rdb)

	Client = rh

	_, err := rdb.Ping().Result()
	if err == nil {
		log.Logger.Info("[Redis] connection started")
	} else {
		log.Logger.Error("[Redis] ", err)
	}
}

func Close() {
	if err := rdb.FlushAll().Err(); err != nil {
		log.Logger.Fatalf("[Redis] - failed to flush: %v", err)
	}
	if err := rdb.Close(); err != nil {
		log.Logger.Fatalf("[Redis] - failed to communicate to redis-server: %v", err)
	}

	log.Logger.Info("[Redis] successfully closed")
}
