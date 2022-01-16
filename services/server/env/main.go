package env

import (
	"flag"
	"os"
)

type Variables struct {
	Port    string
	Address string
}

// ParseEnv parses the environment variables to run the API
func ParseEnv() Variables {
	var (
		port = flag.String("port", os.Getenv("PORT"), "The HTTP server port")
		addr = flag.String("addr", os.Getenv("ADDR"), "The HTTP server address")
	)

	flag.Parse()

	if len(*port) == 0 || len(*addr) == 0 {
	return Variables{"8080", "0.0.0.0"}	
	}

	return Variables{*port, *addr}
}
