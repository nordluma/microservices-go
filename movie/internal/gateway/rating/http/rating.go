package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nordluma/microservices-go/movie/internal/gateway"
	"github.com/nordluma/microservices-go/rating/pkg/model"
)

type GateWay struct {
	addr string
}

func NewGateWay(addr string) *GateWay {
	return &GateWay{addr: addr}
}

func (g *GateWay) GetAggregatedRating(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
) (float64, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/rating", nil)
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
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		g.addr+"/rating",
		nil,
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
