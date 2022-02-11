package env

import (
	"flag"
	"os"
)

type RedisConn struct {
	Address  string
	Password string
	DB       int
}

type Variables struct {
	Port    string
	Address string
	Redis   RedisConn
}

var Vars = Variables{
	Port:    "8080",
	Address: "0.0.0.0",
	Redis: RedisConn{
		Address:  "redis:6379",
		Password: "",
		DB:       0,
	},
}

func Load() {
	var (
		port      = flag.String("port", os.Getenv("PORT"), "The HTTP server port")
		addr      = flag.String("addr", os.Getenv("ADDR"), "The HTTP server address")
		redisAddr = flag.String("redisAddr", os.Getenv("REDIS_ADDR"), "Redis address domain")
		redisPw   = flag.String("redisPw", os.Getenv("REDIS_PW"), "Redis password")
	)

	flag.Parse()

	Vars.Redis = RedisConn{*redisAddr, *redisPw, 0}
	
	if *addr != "" {
		Vars.Address = *addr
	}
	if *port != "" {
		Vars.Port = *port
	}
}
