package api

import (
	"streamserver/httpserver"
	"streamserver/log"
	"streamserver/streaming"
	"streamserver/streaming/hub"
	streamingTransport "streamserver/streaming/transport"
)

func Run() {
	go hub.Instance.Listen()

	s, err := httpserver.New()
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	streamingTransport.NewHTTP(streaming.Initialize(), s.Router)
	streamingTransport.NewWS(streaming.Initialize(), s.Router)

	s.Run()
}
