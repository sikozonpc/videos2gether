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

var envVars = env.ParseEnv()

// Server struct
type Server struct {
	Router *mux.Router
	// Logger *provider.Logger
}

// New creates a new http server instance
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

// Run the server instance
func (s *Server) Run() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // TODO: Correctly setup CORS
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
	})

	port := ":" + envVars.Port

	h := &http.Server{
		Addr:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      c.Handler(s),
	}

	go func() {
		log.Printf("Listening on %s\n", port)

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
