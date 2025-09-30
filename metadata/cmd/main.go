package main

import (
	"log"

	"google.golang.org/grpc/metadata"
)

func main() {
	log.Println("Starting movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
}
