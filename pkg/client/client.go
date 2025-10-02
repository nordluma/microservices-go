package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const contentType = "application/x-www-form-urlencoded"

type Client struct {
	addr   string
	client *http.Client
}

func NewClient(addr string) *Client {
	return &Client{
		addr:   fmt.Sprintf("http://%s", addr),
		client: http.DefaultClient,
	}
}

func (c *Client) Register(
	ctx context.Context,
	instanceId, serviceName, hostPort string,
) error {
	params := url.Values{}
	params.Add("instanceId", instanceId)
	params.Add("serviceName", serviceName)
	params.Add("hostPort", hostPort)
	buf := bytes.NewBufferString(params.Encode())

	req, err := http.NewRequestWithContext(ctx, "POST", c.addr+"/register", buf)
	if err != nil {
		return err
	}
	req.Header.Set("content-type", contentType)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("non ok response: %d", res.StatusCode)
	}

	return nil
}

func (c *Client) Deregister(
	ctx context.Context,
	instanceID, serviceName string,
) error {
	params := url.Values{}
	params.Add("instanceId", instanceID)
	params.Add("serviceName", serviceName)
	buf := bytes.NewBufferString(params.Encode())

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.addr+"/deregister",
		buf,
	)
	if err != nil {
		return err
	}

	req.Header.Set("content-type", contentType)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("non ok response: %d", res.StatusCode)
	}

	return nil
}

func (c *Client) Discover(
	ctx context.Context,
	serviceName string,
) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.addr+"/discover", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("serviceName", serviceName)
	req.URL.RawQuery = query.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non ok response: %d", res.StatusCode)
	}

	var instanceAddrs []string
	if err := json.NewDecoder(res.Body).Decode(&instanceAddrs); err != nil {
		return nil, err
	}

	return instanceAddrs, nil
}

func (c *Client) HealthCheck(
	instanceId, serviceName string,
) error {
	req, err := http.NewRequest("GET", c.addr+"/healthz", nil)
	if err != nil {
		return err
	}

	query := req.URL.Query()
	query.Add("instanceId", instanceId)
	query.Add("serviceName", serviceName)
	req.URL.RawQuery = query.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("non ok response: %d", res.StatusCode)
	}

	return nil
}
