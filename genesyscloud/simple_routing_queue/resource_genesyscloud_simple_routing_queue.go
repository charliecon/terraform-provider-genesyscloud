package simple_routing_queue

import (
	"context"
	"fmt"
	"log"
	"terraform-provider-genesyscloud/genesyscloud/consistency_checker"
	"terraform-provider-genesyscloud/genesyscloud/provider"
	"terraform-provider-genesyscloud/genesyscloud/util"
	"terraform-provider-genesyscloud/genesyscloud/util/constants"
	"time"

	resourceExporter "terraform-provider-genesyscloud/genesyscloud/resource_exporter"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v133/platformclientv2"
)

func getAllRoutingQueues(ctx context.Context, clientConfig *platformclientv2.Configuration) (resourceExporter.ResourceIDMetaMap, diag.Diagnostics) {
	resources := make(resourceExporter.ResourceIDMetaMap)
	proxy := getSimpleRoutingQueueProxy(clientConfig)

	// Newly created resources often aren't returned unless there's a delay
	time.Sleep(5 * time.Second)

	queues, resp, err := proxy.getAllRoutingQueues(ctx)
	if err != nil {
		return nil, util.BuildAPIDiagnosticError(resourceName, fmt.Sprintf("failed to get routing queues: %v", err), resp)
	}

	for _, queue := range *queues {
		resources[*queue.Id] = &resourceExporter.ResourceMeta{Name: *queue.Name}
	}

	return resources, nil
}

/*
The resource_genesyscloud_simple_routing_queue.go contains all of the methods that perform the core logic for a resource.
In general a resource should have a approximately 5 methods in it:

1.  A create.... function that the resource will use to create a Genesys Cloud object (e.g. genesyscloud_simple_routing_queue)
2.  A read.... function that looks up a single resource.
3.  An update... function that updates a single resource.
4.  A delete.... function that deletes a single resource.

Two things to note:

1.  All code in these methods should be focused on getting data in and out of Terraform.  All code that is used for interacting
    with a Genesys API should be encapsulated into a proxy class contained within the package.

2.  In general, to keep this file somewhat manageable, if you find yourself with a number of helper functions move them to a
utils file in the package.  This will keep the code manageable and easy to work through.
*/

// createSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to create a simple queue in Genesys cloud.
func createSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// CREATE-TODO 1: Get an instance of the proxy (example can be found in the delete method below)

	// CREATE-TODO 2: Create variables for each field in our schema.ResourceData object
	// Example:
	name := d.Get("name").(string)

	log.Printf("Creating simple queue %s", name)

	// CREATE-TODO 3: Create a queue struct using the Genesys Cloud platform go sdk

	// CREATE-TODO 4: Call the proxy function to create our queue. The proxy function we want to use here is createRoutingQueue(ctx context.Context, queue *platformclientv2.Createqueuerequest)
	// Note: We won't need the response object returned. Also, don't forget about error handling!

	// CREATE-TODO 5: Set ID in the schema.ResourceData object

	return readSimpleRoutingQueue(ctx, d, meta)
}

// readSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to read a simple queue from Genesys cloud.
func readSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// CREATE-TODO 1: Get an instance of the proxy

	log.Printf("Reading simple queue %s", d.Id())
	cc := consistency_checker.NewConsistencyCheck(ctx, d, meta, ResourceSimpleRoutingQueue(), constants.DefaultConsistencyChecks, resourceName)
	return util.WithRetriesForRead(ctx, d, func() *retry.RetryError {
		/*
			CREATE-TODO 2: Call the proxy function getRoutingQueue(ctx context.Context, id string) to find our queue, passing in the ID from the resource data object
			The returned value are: Queue (*platformclientv2.Queue), Status Code (int), error
			If the error is not nil, we should pass the status code to the function gcloud.IsStatus404ByInt(int)
			If the status code is 404, return a resource.RetryableError. Otherwise, it should be a NonRetryableError
		*/

		// CREATE-TODO 3: Set our values in the schema resource data, based on the values in the Queue object returned from the API

		return cc.CheckState(d)
	})
}

// updateSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to update a simple queue in Genesys cloud.
func updateSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("Updating simple queue %s", d.Id())
	// CREATE-TODO 1: Get an instance of the proxy

	// CREATE-TODO 2: Create variables for each field in our schema.ResourceData object

	// CREATE-TODO	3: Create a queue struct using the Genesys Cloud platform go sdk

	// CREATE-TODO 4: Call the proxy function updateRoutingQueue(context.Context, id string, *platformclientv2.Queuerequest) to update our queue

	return readSimpleRoutingQueue(ctx, d, meta)
}

// deleteSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to delete a simple queue from Genesys cloud.
func deleteSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// CREATE-TODO 1: Get an instance of the proxy (done)
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	proxy := getSimpleRoutingQueueProxy(sdkConfig)

	/*
		CREATE-TODO 2: Call the proxy function deleteRoutingQueue(ctx context.Context, id string)
		Again, we won't be needing the returned response object
	*/

	log.Printf("Deleting simple queue %s", d.Id())
	// Check that queue has been deleted by trying to get it from the API
	return util.WithRetries(ctx, 30*time.Second, func() *retry.RetryError {
		_, resp, err := proxy.getRoutingQueue(ctx, d.Id())

		if err == nil {
			return retry.NonRetryableError(fmt.Errorf("error deleting routing queue %s: %s", d.Id(), err))
		}
		if util.IsStatus404(resp) {
			// Success: Routing Queue deleted
			log.Printf("Deleted routing queue %s", d.Id())
			return nil
		}

		return retry.RetryableError(fmt.Errorf("routing queue %s still exists", d.Id()))
	})
}
