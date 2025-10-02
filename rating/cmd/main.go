package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nordluma/microservices-go/pkg/client"
	"github.com/nordluma/microservices-go/pkg/discovery"
	"github.com/nordluma/microservices-go/rating/internal/controller/rating"
	httpHandler "github.com/nordluma/microservices-go/rating/internal/handler/http"
	"github.com/nordluma/microservices-go/rating/internal/repository/memory"
)

const (
	serviceName = "rating"
	port        = 8082
)

func main() {
	log.Println("starting rating service")
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

	repo := memory.NewRepository()
	ctrl := rating.NewController(repo)
	handler := httpHandler.NewHandler(ctrl)

	http.HandleFunc("/rating", http.HandlerFunc(handler.Handle))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalln(err)
	}
}
