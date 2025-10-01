package main

import (
	"log"
	"net/http"

	"github.com/nordluma/microservices-go/rating/internal/controller/rating"
	httpHandler "github.com/nordluma/microservices-go/rating/internal/handler/http"
	"github.com/nordluma/microservices-go/rating/internal/repository/memory"
)

func main() {
	log.Println("starting rating service")
	repo := memory.NewRepository()
	ctrl := rating.NewController(repo)
	handler := httpHandler.NewHandler(ctrl)

	http.HandleFunc("/rating", http.HandlerFunc(handler.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalln(err)
	}
}
