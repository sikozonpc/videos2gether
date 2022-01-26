package env

import (
	"flag"
	"os"
)

type Variables struct {
	Port    string
	Address string
}

var Vars = Variables{
	Port:    "8080",
	Address: "0.0.0.0",
}

func init() {
	var (
		port = flag.String("port", os.Getenv("PORT"), "The HTTP server port")
		addr = flag.String("addr", os.Getenv("ADDR"), "The HTTP server address")
	)

	flag.Parse()

	if *addr != "" {
		Vars.Address = *addr
	}
	if *port != "" {
		Vars.Port = *port
	}
}
