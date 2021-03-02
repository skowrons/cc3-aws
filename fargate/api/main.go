package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// TODO: implement functionality
func (s *server) routes() {
    log.Println("Register routes.")
    s.router.HandleFunc("/", s.handleRoot())
	s.router.HandleFunc("/_healthcheck", s.handleHealthCheck())
}

func (s *server) handleRoot() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("hello"))
    }
}

// health check for aws fargate
func (s *server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthcheck okay!"))
	}
}

func main() {
    port := ":8080"
	log.Printf("Starting server on port %s\n", port)

	handler := &server{
		router: mux.NewRouter(),
	}
	handler.routes()

	s := &http.Server{
		Addr:         port,
		Handler:      handler,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Fatal(s.ListenAndServe())
}
