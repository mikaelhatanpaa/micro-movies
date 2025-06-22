package memory

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"movieexample.com/pkg/discovery"
)

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

// Registry defines an in-memory service regisry.
// Note: this registry does not perform health monitoring of active instances.
type Registry struct {
	sync.RWMutex
	serviceAddr map[string]map[string]*serviceInstance
}

// NewRegistry creates a new in-memory service registry instance.
func NewRegistry() *Registry {
	return &Registry{serviceAddr: map[string]map[string]*serviceInstance{}}
}

// Register creates a service record in the registry.
func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddr[serviceName]; !ok {
		r.serviceAddr[serviceName] = map[string]*serviceInstance{}
	}

	r.serviceAddr[serviceName][instanceID] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

// Deregister removes a service record from the registry.
func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddr[serviceName]; !ok {
		log.Printf("unknown servicename: %s", serviceName)
		return nil
	}

	delete(r.serviceAddr[serviceName], instanceID)

	return nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddr[serviceName]; !ok {
		return errors.New("instance " + instanceID + " of service " + serviceName + " is not registered yet")
	}
	if _, ok := r.serviceAddr[serviceName][instanceID]; !ok {
		return errors.New("service instance is not registered yet")
	}
	r.serviceAddr[serviceName][instanceID].lastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.serviceAddr[serviceName]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for instanceID, i := range r.serviceAddr[serviceName] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			log.Println("Instance " + instanceID + " of service " + serviceName + " is not active, skipping")
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}
