package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nordluma/microservices-go/metadata/internal/controller/metadata"
	httpHandler "github.com/nordluma/microservices-go/metadata/internal/handler/http"
	"github.com/nordluma/microservices-go/metadata/internal/repository/memory"
	"github.com/nordluma/microservices-go/pkg/client"
	"github.com/nordluma/microservices-go/pkg/discovery"
)

const (
	serviceName = "metadata"
	port        = 8081
)

func main() {
	log.Println("Starting movie metadata service")
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

	repo := memory.New()
	controller := metadata.New(repo)
	handler := httpHandler.New(controller)

	http.Handle("/metadata", http.HandlerFunc(handler.GetMetadata))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}
