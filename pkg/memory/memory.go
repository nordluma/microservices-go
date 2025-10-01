package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/nordluma/microservices-go/pkg/discovery"
)

type (
	serviceNameType string
	instanceIDType  string
)

type InMemoryRegistry struct {
	serviceAddrs map[serviceNameType]map[instanceIDType]*serviceInstance

	sync.RWMutex
}

type serviceInstance struct {
	hostPort      string
	lastHeartbeat time.Time
}

func NewRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{
		serviceAddrs: make(
			map[serviceNameType]map[instanceIDType]*serviceInstance,
		),
	}
}

func (r *InMemoryRegistry) Register(
	ctx context.Context,
	instanceID, serviceName, hostPort string,
) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		r.serviceAddrs[serviceNameType(serviceName)] = map[instanceIDType]*serviceInstance{}
	}

	instance := &serviceInstance{
		hostPort:      hostPort,
		lastHeartbeat: time.Now(),
	}
	r.serviceAddrs[serviceNameType(serviceName)][instanceIDType(instanceID)] = instance

	return nil
}

func (r *InMemoryRegistry) Deregister(
	ctx context.Context,
	instancedID, serviceName string,
) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		return discovery.ErrNotFound
	}

	delete(
		r.serviceAddrs[serviceNameType(serviceName)],
		instanceIDType(instancedID),
	)

	return nil
}

func (r *InMemoryRegistry) HealthCheck(instanceID, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		return errors.New("service not registered")
	}

	if _, ok := r.serviceAddrs[serviceNameType(serviceName)][instanceIDType(instanceID)]; !ok {
		return errors.New("service instance not registered")
	}

	r.serviceAddrs[serviceNameType(serviceName)][instanceIDType(instanceID)].lastHeartbeat = time.Now()

	return nil
}

func (r *InMemoryRegistry) Discover(
	ctx context.Context,
	serviceName string,
) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.serviceAddrs[serviceNameType(serviceName)]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string
	for _, v := range r.serviceAddrs[serviceNameType(serviceName)] {
		if time.Since(v.lastHeartbeat) > 5*time.Second {
			continue
		}

		res = append(res, v.hostPort)
	}

	return res, nil
}
