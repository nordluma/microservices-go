package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/nordluma/microservices-go/movie/internal/gateway"
	"github.com/nordluma/microservices-go/pkg/discovery"
	"github.com/nordluma/microservices-go/rating/pkg/model"
)

type GateWay struct {
	registry discovery.Registry
}

func NewGateWay(registry discovery.Registry) *GateWay {
	return &GateWay{registry: registry}
}

func (g *GateWay) GetAggregatedRating(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
) (float64, error) {
	url, err := g.getServiceUrl(ctx)
	if err != nil {
		return 0, err
	}

	log.Printf("calling rating service. Request: GET %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	queryParams := req.URL.Query()
	queryParams.Add("id", string(recordID))
	queryParams.Add("type", fmt.Sprintf("%v", recordType))
	req.URL.RawQuery = queryParams.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if res.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non ok response: %d", res.StatusCode)
	}

	var rating float64
	if err := json.NewDecoder(res.Body).Decode(&rating); err != nil {
		return 0, err
	}

	return rating, nil
}

func (g *GateWay) InsertRating(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
	userID model.UserID,
	value model.RatingValue,
) error {
	url, err := g.getServiceUrl(ctx)
	if err != nil {
		return err
	}

	log.Printf("calling rating service. Request: %s", url)
	req, err := http.NewRequestWithContext(
		ctx, http.MethodPut, url, nil,
	)
	if err != nil {
		return err
	}

	queryParams := req.URL.Query()
	queryParams.Add("id", string(recordID))
	queryParams.Add("type", string(recordType))
	queryParams.Add("userId", string(userID))
	queryParams.Add("value", fmt.Sprintf("%d", value))
	req.URL.RawQuery = queryParams.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("non ok response: %d", res.StatusCode)
	}

	return nil
}

func (g *GateWay) getServiceUrl(ctx context.Context) (string, error) {
	addrs, err := g.registry.Discover(ctx, "rating")
	if err != nil {
		return "", err
	}

	if len(addrs) == 0 {
		return "", fmt.Errorf("no rating service instances available")
	}

	return "http://" + addrs[rand.Intn(len(addrs))] + "/rating", nil
}
