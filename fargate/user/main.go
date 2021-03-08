package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) routes() {
	log.Println("Register routes.")
	s.router.HandleFunc("/", healthCheckMiddleware(s.handleRoot())) // only used for standard healthcheck
	s.router.HandleFunc("/users", s.handleGetUsers()).Methods(http.MethodGet)
}

// healthcheck for aws fargate -> standard path is root
func healthCheckMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.UserAgent(), "ELB-HealthChecker") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("healthcheck okay!"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *server) handleGetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := []struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{
			{"Tom", 32},
			{"Jim", 55},
			{"Luna", 19}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}
}

func main() {
	port := ":8080"
	log.Printf("Starting user server on port %s\n", port)

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
