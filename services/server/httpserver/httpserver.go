package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"streamserver/env"
	"time"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
)

type Server struct {
	Router *mux.Router
}

func New() (*Server, error) {
	r := mux.NewRouter()
	s := &Server{
		Router: r,
	}

	s.Router.HandleFunc("/health", handleHealthCheck)

	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) Run() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{env.Vars.Cors.Origin},
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
	})

	fullAddr := env.Vars.Address + ":" + env.Vars.Port

	h := &http.Server{
		Addr:         fullAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      c.Handler(s),
	}

	go func() {
		log.Printf("Listening on %s\n", fullAddr)

		if err := h.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Println("\nGracefully shutting down the server...")
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All good to Go :)")
}
