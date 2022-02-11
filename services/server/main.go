package main

import (
	"streamserver/cmd/api"
	"streamserver/env"
	"streamserver/redis"
)
 
func main() {
	env.Load()

	redis.Connect()
	defer redis.Close()

	api.Run()
}
