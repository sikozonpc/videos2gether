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

type Cors struct {
	Origin string
}

type Variables struct {
	Port    string
	Address string
	APIKey  string
	Redis   RedisConn
	Cors    Cors
}

var Vars = Variables{
	Port:    "8080",
	Address: "0.0.0.0",
	Redis: RedisConn{
		Address:  "redis:6379",
		Password: "",
		DB:       0,
	},
	APIKey: "",
	Cors: Cors{
		Origin: "*",
	},
}

func Load() {
	var (
		port       = flag.String("port", os.Getenv("PORT"), "The HTTP server port")
		addr       = flag.String("addr", os.Getenv("ADDR"), "The HTTP server address")
		redisAddr  = flag.String("redisAddr", os.Getenv("REDIS_ADDR"), "Redis address domain")
		redisPw    = flag.String("redisPw", os.Getenv("REDIS_PW"), "Redis password")
		apiKey     = flag.String("apiKey", os.Getenv("API_KEY"), "Auth API key")
		corsOrigin = flag.String("corsOrigin", os.Getenv("CORS_ORIGIN"), "Cors origin")
	)

	flag.Parse()

	Vars.Redis = RedisConn{*redisAddr, *redisPw, 0}
	Vars.APIKey = *apiKey

	if *corsOrigin != "" {
		Vars.Cors = Cors{Origin: *corsOrigin}
	}

	if *addr != "" {
		Vars.Address = *addr
	}
	if *port != "" {
		Vars.Port = *port
	}
}
