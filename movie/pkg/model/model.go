package model

import "github.com/nordluma/microservices-go/metadata/pkg/model"

type MovieDetails struct {
	Rating   *float64       `json:"rating"`
	Metadata model.Metadata `json:"metadata"`
}
