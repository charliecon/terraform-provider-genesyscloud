package simple_routing_queue

import (
	"context"
	"fmt"
	rc "terraform-provider-genesyscloud/genesyscloud/resource_cache"

	"github.com/mypurecloud/platform-client-sdk-go/v133/platformclientv2"
)

type createSimpleRoutingQueueFunc func(context.Context, *simpleRoutingQueueProxy, *platformclientv2.Createqueuerequest) (*platformclientv2.Queue, *platformclientv2.APIResponse, error)
type getSimpleRoutingQueueByIdFunc func(context.Context, *simpleRoutingQueueProxy, string) (*platformclientv2.Queue, *platformclientv2.APIResponse, error)
type getAllSimpleRoutingQueuesFunc func(context.Context, *simpleRoutingQueueProxy) (*[]platformclientv2.Queue, *platformclientv2.APIResponse, error)
type updateSimpleRoutingQueueFunc func(context.Context, *simpleRoutingQueueProxy, string, *platformclientv2.Queuerequest) (*platformclientv2.Queue, *platformclientv2.APIResponse, error)
type deleteSimpleRoutingQueueFunc func(context.Context, *simpleRoutingQueueProxy, string) (*platformclientv2.APIResponse, error)
type getSimpleRoutingQueueIdByNameFunc func(context.Context, *simpleRoutingQueueProxy, string) (id string, response *platformclientv2.APIResponse, err error, retryable bool)

var internalProxy *simpleRoutingQueueProxy

type simpleRoutingQueueProxy struct {
	routingApi                        *platformclientv2.RoutingApi
	createSimpleRoutingQueueAttr      createSimpleRoutingQueueFunc
	getSimpleRoutingQueueByIdAttr     getSimpleRoutingQueueByIdFunc
	getAllSimpleRoutingQueuesAttr     getAllSimpleRoutingQueuesFunc
	getSimpleRoutingQueueIdByNameAttr getSimpleRoutingQueueIdByNameFunc
	updateSimpleRoutingQueueAttr      updateSimpleRoutingQueueFunc
	deleteSimpleRoutingQueueAttr      deleteSimpleRoutingQueueFunc
	simpleRoutingQueueCache           rc.CacheInterface[platformclientv2.Queue]
}

// newSimpleRoutingQueueProxy initializes the simple routing queue proxy with all the data needed to communicate with Genesys Cloud
func newSimpleRoutingQueueProxy(clientConfig *platformclientv2.Configuration) *simpleRoutingQueueProxy {
	api := platformclientv2.NewRoutingApiWithConfig(clientConfig)
	simpleRoutingQueueCache := rc.NewResourceCache[platformclientv2.Queue]()

	return &simpleRoutingQueueProxy{
		routingApi:                        api,
		createSimpleRoutingQueueAttr:      createSimpleRoutingQueueFn,
		getSimpleRoutingQueueByIdAttr:     getSimpleRoutingQueueByIdFn,
		getSimpleRoutingQueueIdByNameAttr: getSimpleRoutingQueueIdByNameFn,
		updateSimpleRoutingQueueAttr:      updateSimpleRoutingQueueFn,
		deleteSimpleRoutingQueueAttr:      deleteSimpleRoutingQueueFn,
		simpleRoutingQueueCache:           simpleRoutingQueueCache,
	}
}

// getSimpleRoutingQueueProxy acts as a singleton to for the internalProxy.  It also ensures
// that we can still proxy our tests by directly setting internalProxy package variable
func getSimpleRoutingQueueProxy(clientConfig *platformclientv2.Configuration) *simpleRoutingQueueProxy {
	if internalProxy == nil {
		internalProxy = newSimpleRoutingQueueProxy(clientConfig)
	}
	return internalProxy
}

// createRoutingQueue creates a Genesys Cloud Routing Queue
func (p *simpleRoutingQueueProxy) createSimpleRoutingQueue(ctx context.Context, queue *platformclientv2.Createqueuerequest) (*platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return p.createSimpleRoutingQueueAttr(ctx, p, queue)
}

// getRoutingQueue retrieves a Genesys Cloud Routing Queue by ID
func (p *simpleRoutingQueueProxy) getSimpleRoutingQueue(ctx context.Context, id string) (*platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return p.getSimpleRoutingQueueByIdAttr(ctx, p, id)
}

func (p *simpleRoutingQueueProxy) getAllSimpleRoutingQueues(ctx context.Context) (*[]platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return p.getAllSimpleRoutingQueuesAttr(ctx, p)
}

// getRoutingQueueIdByName retrieves a Genesys Cloud Routing Queue ID by its name
func (p *simpleRoutingQueueProxy) getSimpleRoutingQueueIdByName(ctx context.Context, name string) (string, *platformclientv2.APIResponse, error, bool) {
	return p.getSimpleRoutingQueueIdByNameAttr(ctx, p, name)
}

// updateRoutingQueue updates a Genesys Cloud Routing Queue
func (p *simpleRoutingQueueProxy) updateSimpleRoutingQueue(ctx context.Context, id string, queue *platformclientv2.Queuerequest) (*platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return p.updateSimpleRoutingQueueAttr(ctx, p, id, queue)
}

// deleteRoutingQueue deletes a Genesys Cloud Routing Queue
func (p *simpleRoutingQueueProxy) deleteSimpleRoutingQueue(ctx context.Context, id string) (*platformclientv2.APIResponse, error) {
	return p.deleteSimpleRoutingQueueAttr(ctx, p, id)
}

// createRoutingQueueFn is an implementation function for creating a Genesys Cloud Routing Queue
func createSimpleRoutingQueueFn(ctx context.Context, proxy *simpleRoutingQueueProxy, queue *platformclientv2.Createqueuerequest) (*platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return proxy.routingApi.PostRoutingQueues(*queue)
}

func getSimpleRoutingQueueByIdFn(ctx context.Context, proxy *simpleRoutingQueueProxy, id string) (*platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return proxy.routingApi.GetRoutingQueue(id)
}

func getAllSimpleRoutingQueuesFn(ctx, proxy *simpleRoutingQueueProxy) (*[]platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	const pageSize = 100
	var allQueues []platformclientv2.Queue

	queues, resp, err := proxy.routingApi.GetRoutingQueues(1, pageSize, "", "", nil, nil, nil, "", false)
	if err != nil {
		return nil, resp, err
	}

	if queues.Entities == nil || len(*queues.Entities) == 0 {
		return &allQueues, nil, nil
	}

	allQueues = append(allQueues, *queues.Entities...)

	for pageNum := 2; pageNum <= *queues.PageCount; pageNum++ {
		queues, resp, err := proxy.routingApi.GetRoutingQueues(pageNum, pageSize, "", "", nil, nil, nil, "", false)
		if err != nil {
			return &allQueues, resp, err
		}

		if queues.Entities == nil || len(*queues.Entities) == 0 {
			break
		}

		allQueues = append(allQueues, *queues.Entities...)
	}

	for _, queue := range allQueues {
		rc.SetCache(proxy.simpleRoutingQueueCache, *queue.Id, queue)
	}

	return &allQueues, nil, nil
}

func getSimpleRoutingQueueIdByNameFn(ctx context.Context, proxy *simpleRoutingQueueProxy, name string) (string, *platformclientv2.APIResponse, error, bool) {
	const pageSize = 100
	notFoundError := fmt.Errorf("no routing queues found with name %s", name)

	queues, resp, getErr := proxy.routingApi.GetRoutingQueues(1, pageSize, "", name, nil, nil, nil, "", false)
	if getErr != nil {
		return "", resp, getErr, false
	}

	if queues.Entities == nil || len(*queues.Entities) == 0 {
		return "", nil, notFoundError, true
	}

	for _, queue := range *queues.Entities {
		if queue.Name != nil && *queue.Name == name {
			return *queue.Id, nil, nil, false
		}
	}

	for pageNum := 2; pageNum <= *queues.PageCount; pageNum++ {
		queues, resp, getErr := proxy.routingApi.GetRoutingQueues(pageNum, pageSize, "", name, nil, nil, nil, "", false)
		if getErr != nil {
			return "", resp, getErr, false
		}

		if queues.Entities == nil || len(*queues.Entities) == 0 {
			return "", nil, notFoundError, true
		}

		for _, queue := range *queues.Entities {
			if queue.Name != nil && *queue.Name == name {
				return *queue.Id, nil, nil, false
			}
		}
	}

	return "", nil, notFoundError, true
}

func updateSimpleRoutingQueueFn(_ context.Context, proxy *simpleRoutingQueueProxy, id string, body *platformclientv2.Queuerequest) (*platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return proxy.routingApi.PutRoutingQueue(id, *body)
}

func deleteSimpleRoutingQueueFn(_ context.Context, proxy *simpleRoutingQueueProxy, id string) (*platformclientv2.APIResponse, error) {
	return proxy.routingApi.DeleteRoutingQueue(id, true)
}
