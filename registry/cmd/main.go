package main

import (
	"fmt"
	"log"
	"net/http"

	inmemory "github.com/nordluma/microservices-go/pkg/inMemory"
	"github.com/nordluma/microservices-go/registry/internal/handler"
)

const port = 8080

func main() {
	log.Println("starting registry service")
	register := inmemory.NewRegistry()
	handler := handler.NewHandler(register)

	http.HandleFunc("/register", http.HandlerFunc(handler.Register))
	http.HandleFunc("/deregister", http.HandlerFunc(handler.Deregister))
	http.HandleFunc("/discover", http.HandlerFunc(handler.Discover))
	http.HandleFunc("/healthz", http.HandlerFunc(handler.HealthCheck))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalln(err)
	}
}
