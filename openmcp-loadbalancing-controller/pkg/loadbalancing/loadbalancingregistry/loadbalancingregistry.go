// Package registry defines the Registry interface which can be used with goproxy.
package loadbalancingregistry

import (
	"errors"
	"log"
	"openmcp/openmcp/omcplog"
	"sync"
)

var lock sync.RWMutex

// Common errors.
var (
	ErrServiceNotFound = errors.New("service name/version not found")
)

// Registry is an interface used to lookup the target host
// for a given service name / version pair.
type Registry interface {
	Add(host, path, endpoint string)                // Add an endpoint to our registry
	Delete(host, path, endpoint string)             // Remove an endpoint to our registry
	Failure(host, path, endpoint string, err error) // Mark an endpoint as failed.
	Lookup(host, path string) (string, error)       // Return the endpoint list for the given service name/version
	IngressDelete(host, path string)
	//IngressLookup(host string, path string, endpoint string) (bool)
}

// DefaultRegistry is a basic registry using the following format:
// {
//   "Host": {
//     "Path": [
//       "cluster1",
//       "cluster2"
//     ],
//   },
// }

// DefaultRegistry is a basic registry using the following format:
// {
//   "Host": {
//     "Path": [
//       "serviceName",
//     ],
//   },
// }

//type DefaultRegistry map[string]map[string]map[string]string
type DefaultRegistry map[string]map[string]string

// Lookup return the endpoint list for the given service name/version.

func (r DefaultRegistry) Lookup(host string, path string) (string, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(LoadbalancingRegistry)] Function Lookup")
	//fmt.Println("----Lookup----")
	lock.RLock()
	target, ok := r[host][path]
	lock.RUnlock()
	if !ok {
		return "", ErrServiceNotFound
	}
	return target, nil
}


func (r DefaultRegistry) Failure(host, path, endpoint string, err error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(LoadbalancingRegistry)] Function Failure")
	// Would be used to remove an endpoint from the rotation, log the failure, etc.
	//log.Printf("Error accessing %s/%s (%s): %s", path, endpoint, err)
	log.Printf("Error accessing %s %s (%s): %s", host, path, endpoint, err)
}

func (r DefaultRegistry) Add(host, path, endpoint string) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(LoadbalancingRegistry)] Function Add")
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[host]
	if !ok {
		service = map[string]string{}
		r[host] = service
	}
	service[path] = endpoint
	//service[path] = append(service[path], endpoint)
}

// Delete removes the given endpoit for the service name/version.
func (r DefaultRegistry) Delete(host, path, endpoint string) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(LoadbalancingRegistry)] Function Delete")
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[host]
	if !ok {
		return
	}
	omcplog.V(5).Info(service)
}

//// Delete removes the given endpoit for the service name/version.
//func (r DefaultRegistry) Delete(host, path, endpoint string) {
//	fmt.Println("----Delete----")
//	lock.Lock()
//	defer lock.Unlock()
//
//	service, ok := r[host]
//	if !ok {
//		return
//	}
//
//begin:
//	for i, svc := range service[path] {
//		if svc == endpoint {
//			copy(service[path][i:], service[path][i+1:])
//			service[path] = service[path][:len(service[path])-1]
//			goto begin
//		}
//	}
//	fmt.Println("Delete test")
//}

func (r DefaultRegistry) IngressDelete(host, path string) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(LoadbalancingRegistry)] Function IngressDelete")
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[host]
	if !ok {
		return
	}
	delete(service, path)

	if len(r[host]) == 0 {
		delete(r, host)
	}
}
