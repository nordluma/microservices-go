package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nordluma/microservices-go/metadata/pkg/model"
	"github.com/nordluma/microservices-go/movie/internal/gateway"
)

type Gateway struct {
	addr string
}

func NewGateway(addr string) *Gateway {
	return &Gateway{addr: addr}
}

func (g *Gateway) GetById(
	ctx context.Context,
	id string,
) (*model.Metadata, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/metadata", nil)
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
