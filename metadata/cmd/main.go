package main

import (
	"log"
	"net/http"

	"github.com/nordluma/microservices-go/metadata/internal/controller/metadata"
	httpHandler "github.com/nordluma/microservices-go/metadata/internal/handler/http"
	"github.com/nordluma/microservices-go/metadata/internal/repository/memory"
)

func main() {
	log.Println("Starting movie metadata service")
	repo := memory.New()
	controller := metadata.New(repo)
	handler := httpHandler.New(controller)

	http.Handle("/metadata", http.HandlerFunc(handler.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
