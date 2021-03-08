package main

import (
	"io/ioutil"
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
	s.router.HandleFunc("/", healthCheckMiddleware(s.handleCats()))
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

func (s *server) handleCats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://cat-fact.herokuapp.com/facts")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		
		log.Println(body)

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
