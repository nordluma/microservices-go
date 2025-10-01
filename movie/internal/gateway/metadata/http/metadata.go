package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/nordluma/microservices-go/metadata/pkg/model"
	"github.com/nordluma/microservices-go/movie/internal/gateway"
	"github.com/nordluma/microservices-go/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) GetById(
	ctx context.Context,
	id string,
) (*model.Metadata, error) {
	url, err := g.getServiceUrl(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Calling metadata service. Request: %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	queryParams := req.URL.Query()
	queryParams.Add("id", id)
	req.URL.RawQuery = queryParams.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if res.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non ok response: %d", res.StatusCode)
	}

	var data *model.Metadata
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func (g *Gateway) getServiceUrl(ctx context.Context) (string, error) {
	addrs, err := g.registry.Discover(ctx, "metadata")
	if err != nil {
		return "", err
	}

	if len(addrs) == 0 {
		return "", errors.New("no available metadata services")
	}

	return "http://" + addrs[rand.Intn(len(addrs))] + "/metadata", nil
}
