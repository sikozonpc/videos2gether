package api

import (
	"streamserver/httpserver"
	"streamserver/log"
	"streamserver/streaming"
	"streamserver/streaming/hub"
	streamingTransport "streamserver/streaming/transport"
)

// Run the http server
func Run() {
	go hub.Instance.Run()

	s, err := httpserver.New()
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	streamingTransport.NewHTTP(streaming.Initialize(), s.Router)
	streamingTransport.NewWS(streaming.Initialize(), s.Router)

	s.Run()
}
