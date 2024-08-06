package simple_routing_queue

import (
	"context"
	"fmt"
	rc "terraform-provider-genesyscloud/genesyscloud/resource_cache"

	"github.com/mypurecloud/platform-client-sdk-go/v133/platformclientv2"
)

type createRoutingQueueFunc func(context.Context, *simpleRoutingQueueProxy, *platformclientv2.Createqueuerequest) (*platformclientv2.Queue, *platformclientv2.APIResponse, error)
type getRoutingQueueFunc func(context.Context, *simpleRoutingQueueProxy, string) (*platformclientv2.Queue, *platformclientv2.APIResponse, error)
type getAllRoutingQueuesFunc func(context.Context, *simpleRoutingQueueProxy) (*[]platformclientv2.Queue, *platformclientv2.APIResponse, error)
type updateRoutingQueueFunc func(context.Context, *simpleRoutingQueueProxy, string, *platformclientv2.Queuerequest) (*platformclientv2.Queue, *platformclientv2.APIResponse, error)
type deleteRoutingQueueFunc func(context.Context, *simpleRoutingQueueProxy, string, bool) (*platformclientv2.APIResponse, error)
type getRoutingQueueIdByNameFunc func(context.Context, *simpleRoutingQueueProxy, string) (id string, retryable bool, err error)

var internalProxy *simpleRoutingQueueProxy

type simpleRoutingQueueProxy struct {
	routingApi                  *platformclientv2.RoutingApi
	createRoutingQueueAttr      createRoutingQueueFunc
	getRoutingQueueAttr         getRoutingQueueFunc
	getAllRoutingQueuesAttr     getAllRoutingQueuesFunc
	getRoutingQueueIdByNameAttr getRoutingQueueIdByNameFunc
	updateRoutingQueueAttr      updateRoutingQueueFunc
	deleteRoutingQueueAttr      deleteRoutingQueueFunc
	routingQueueCache           rc.CacheInterface[platformclientv2.Queue]
}

// newSimpleRoutingQueueProxy initializes the simple routing queue proxy with all the data needed to communicate with Genesys Cloud
func newSimpleRoutingQueueProxy(clientConfig *platformclientv2.Configuration) *simpleRoutingQueueProxy {
	api := platformclientv2.NewRoutingApiWithConfig(clientConfig)
	routingQueueCache := rc.NewResourceCache[platformclientv2.Queue]()

	return &simpleRoutingQueueProxy{
		routingApi:                  api,
		createRoutingQueueAttr:      createRoutingQueueFn,
		getRoutingQueueAttr:         getRoutingQueueFn,
		getRoutingQueueIdByNameAttr: getRoutingQueueIdByNameFn,
		updateRoutingQueueAttr:      updateRoutingQueueFn,
		deleteRoutingQueueAttr:      deleteRoutingQueueFn,
		routingQueueCache:           routingQueueCache,
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
func (p *simpleRoutingQueueProxy) createRoutingQueue(ctx context.Context, queue *platformclientv2.Createqueuerequest) (*platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return p.createRoutingQueueAttr(ctx, p, queue)
}

// getRoutingQueue retrieves a Genesys Cloud Routing Queue by ID
func (p *simpleRoutingQueueProxy) getRoutingQueue(ctx context.Context, id string) (*platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return p.getRoutingQueueAttr(ctx, p, id)
}

func (p *simpleRoutingQueueProxy) getAllRoutingQueues(ctx context.Context) (*[]platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return p.getAllRoutingQueuesAttr(ctx, p)
}

// getRoutingQueueIdByName retrieves a Genesys Cloud Routing Queue ID by its name
func (p *simpleRoutingQueueProxy) getRoutingQueueIdByName(ctx context.Context, name string) (string, bool, error) {
	return p.getRoutingQueueIdByNameAttr(ctx, p, name)
}

// updateRoutingQueue updates a Genesys Cloud Routing Queue
func (p *simpleRoutingQueueProxy) updateRoutingQueue(ctx context.Context, id string, queue *platformclientv2.Queuerequest) (*platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return p.updateRoutingQueueAttr(ctx, p, id, queue)
}

// deleteRoutingQueue deletes a Genesys Cloud Routing Queue
func (p *simpleRoutingQueueProxy) deleteRoutingQueue(ctx context.Context, id string, forceDelete bool) (*platformclientv2.APIResponse, error) {
	return p.deleteRoutingQueueAttr(ctx, p, id, forceDelete)
}

// createRoutingQueueFn is an implementation function for creating a Genesys Cloud Routing Queue
func createRoutingQueueFn(ctx context.Context, proxy *simpleRoutingQueueProxy, queue *platformclientv2.Createqueuerequest) (*platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return proxy.routingApi.PostRoutingQueues(*queue)
}

func getRoutingQueueFn(ctx context.Context, proxy *simpleRoutingQueueProxy, id string) (*platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	return proxy.routingApi.GetRoutingQueue(id)
}

func getAllRoutingQueuesFn(ctx, proxy *simpleRoutingQueueProxy) (*[]platformclientv2.Queue, *platformclientv2.APIResponse, error) {
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
		rc.SetCache(proxy.routingQueueCache, *queue.Id, queue)
	}

	return &allQueues, nil, nil
}

func getRoutingQueueIdByNameFn(ctx context.Context, proxy *simpleRoutingQueueProxy, name string) (string, bool, error) {
	const pageSize = 100
	notFoundError := fmt.Errorf("no routing queues found with name %s", name)

	queues, _, getErr := proxy.routingApi.GetRoutingQueues(1, pageSize, "", name, nil, nil, nil, "", false)
	if getErr != nil {
		return "", false, getErr
	}

	if queues.Entities == nil || len(*queues.Entities) == 0 {
		return "", true, notFoundError
	}

	for _, queue := range *queues.Entities {
		if queue.Name != nil && *queue.Name == name {
			return *queue.Id, false, nil
		}
	}

	for pageNum := 2; pageNum <= *queues.PageCount; pageNum++ {
		queues, _, getErr := proxy.routingApi.GetRoutingQueues(pageNum, pageSize, "", name, nil, nil, nil, "", false)
		if getErr != nil {
			return "", false, getErr
		}

		if queues.Entities == nil || len(*queues.Entities) == 0 {
			return "", true, notFoundError
		}

		for _, queue := range *queues.Entities {
			if queue.Name != nil && *queue.Name == name {
				return *queue.Id, false, nil
			}
		}
	}

	return "", true, notFoundError
}

func updateRoutingQueueFn(ctx context.Context, proxy *simpleRoutingQueueProxy, id string, body *platformclientv2.Queuerequest) (*platformclientv2.Queue, *platformclientv2.APIResponse, error) {
	queue, response, err := proxy.routingApi.PutRoutingQueue(id, *body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update queue %s: %v", id, err)
	}
	return queue, response, err
}

func deleteRoutingQueueFn(ctx context.Context, proxy *simpleRoutingQueueProxy, id string, forceDelete bool) (*platformclientv2.APIResponse, error) {
	response, err := proxy.routingApi.DeleteRoutingQueue(id, forceDelete)
	if err != nil {
		return nil, fmt.Errorf("failed to delete queue '%s': %v", id, err)
	}
	return response, err
}
