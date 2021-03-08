package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
    s.router.HandleFunc("/", healthCheckMiddleware(s.handleRoot()))
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

func (s *server) handleRoot() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        endpoint := fmt.Sprintf("http://cat.%s:8080/", os.Getenv("COPILOT_SERVICE_DISCOVERY_ENDPOINT"))

        resp, err := http.Get(endpoint)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer resp.Body.Close()
        body, _ := ioutil.ReadAll(r.Body)

        w.WriteHeader(http.StatusOK)
        w.Write(body)
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
