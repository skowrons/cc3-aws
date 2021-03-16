package main

import (
	"log"
	"net/http"
	"strings"
	"time"
	"encoding/json"

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
    s.router.HandleFunc("/", healthCheckMiddleware(s.handleRoot()))
    s.router.HandleFunc("/products", s.handleProduct()).Methods(http.MethodGet)
}

func healthCheckMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)  {
		if strings.Contains(r.UserAgent(), "ELB-HealthChecker") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("healthcheck okay!"))
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func (s *server) handleProduct() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		products := []struct{
			Name string `json:"name"`
			Store int `json:"store"`
			Description string `json:"description"`
		}{
			{"Auto", 20, "Ganz schnell"},
			{"Computer", 42, "Ganz schnell"},
			{"KVM Switch", 82, "Ganz schnell"},
		}

        w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(products)
    }
}

func (s *server) handleRoot() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }
}

func main() {
    port := ":8080"
	log.Printf("Starting product server on port %s\n", port)

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
