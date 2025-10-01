package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var ErrNotFound = errors.New("no service addresses found")

type Registry interface {
	// Register as service with the registry
	Register(
		ctx context.Context,
		instanceID, serviceName, hostPort string,
	) error

	// Deregister a service from the registry
	Deregister(ctx context.Context, instanceID, serviceName string) error

	// Discover a service from the registry
	Discover(ctx context.Context, serviceName string) ([]string, error)

	// Perform a healthcheck for a service instance
	HealthCheck(instanceID, serviceName string) error
}

func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf(
		"%s-%d",
		serviceName,
		rand.New(rand.NewSource(time.Now().UnixNano())).Int(),
	)
}
