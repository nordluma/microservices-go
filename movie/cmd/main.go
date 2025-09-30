package main

import (
	"log"
	"net/http"

	"github.com/nordluma/microservices-go/movie/internal/controller/movie"
	metadataGateway "github.com/nordluma/microservices-go/movie/internal/gateway/metadata/http"
	ratingGateway "github.com/nordluma/microservices-go/movie/internal/gateway/rating/http"
	httpHandler "github.com/nordluma/microservices-go/movie/internal/handler/http"
)

func main() {
	log.Println("Starting movie service")
	metadataGateway := metadataGateway.NewGateway("localhost:8081")
	ratingGateway := ratingGateway.NewGateWay("localhost:8082")

	controller := movie.NewController(ratingGateway, metadataGateway)
	handler := httpHandler.NewHandler(controller)

	http.Handle("/movie", http.HandlerFunc(handler.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatal(err)
	}
}
