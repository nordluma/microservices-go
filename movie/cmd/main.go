package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nordluma/microservices-go/movie/internal/controller/movie"
	metadataGateway "github.com/nordluma/microservices-go/movie/internal/gateway/metadata/http"
	ratingGateway "github.com/nordluma/microservices-go/movie/internal/gateway/rating/http"
	httpHandler "github.com/nordluma/microservices-go/movie/internal/handler/http"
	"github.com/nordluma/microservices-go/pkg/client"
	"github.com/nordluma/microservices-go/pkg/discovery"
)

const (
	serviceName = "movie"
	port        = 8083
)

func main() {
	log.Println("Starting movie service")
	addr := fmt.Sprintf("localhost:%d", port)

	registryClient := client.NewClient("localhost:8080")
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	log.Printf("instance id: %s", instanceID)
	if err := registryClient.Register(ctx, instanceID, serviceName, addr); err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			if err := registryClient.HealthCheck(instanceID, serviceName); err != nil {
				log.Printf("failed to report healthy state: %v", err)
			}

			time.Sleep(1 * time.Second)
		}
	}()
	defer func() {
		if err := registryClient.Deregister(ctx, instanceID, serviceName); err != nil {
			log.Printf("failed to deregister: %v", err)
		}
	}()

	metadataGateway := metadataGateway.NewGateway(registryClient)
	ratingGateway := ratingGateway.NewGateWay(registryClient)

	controller := movie.NewController(ratingGateway, metadataGateway)
	handler := httpHandler.NewHandler(controller)

	http.Handle("/movie", http.HandlerFunc(handler.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}
