package simple_routing_queue

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v105/platformclientv2"
	"log"
	gcloud "terraform-provider-genesyscloud/genesyscloud"
	"terraform-provider-genesyscloud/genesyscloud/consistency_checker"
	"terraform-provider-genesyscloud/genesyscloud/util/resourcedata"
	"time"
)

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
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := getSimpleRoutingQueueProxy(sdkConfig)

	// CREATE-TODO 2: Create variables for each field in our schema.ResourceData object
	name := d.Get("name").(string)
	callingPartyName := d.Get("calling_party_name").(string)
	enableTranscription := d.Get("enable_transcription").(bool)

	log.Printf("Creating simple queue %s", name)

	// CREATE-TODO 3: Create a queue struct using the Genesys Cloud platform go sdk
	queueCreate := &platformclientv2.Createqueuerequest{
		Name:                &name,
		CallingPartyName:    &callingPartyName,
		EnableTranscription: &enableTranscription,
	}

	// CREATE-TODO 4: Call the proxy function to create our queue. The proxy function we want to use here is createRoutingQueue(ctx context.Context, queue *platformclientv2.Createqueuerequest)
	// Note: We won't need the response object returned. Also, don't forget about error handling!
	queueResp, _, err := proxy.createRoutingQueue(ctx, queueCreate)
	if err != nil {
		return diag.Errorf("failed to create queue %s: %v", name, err)
	}

	// CREATE-TODO 5: Set ID in the schema.ResourceData object
	d.SetId(*queueResp.Id)

	return readSimpleRoutingQueue(ctx, d, meta)
}

// readSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to read a simple queue from Genesys cloud.
func readSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// CREATE-TODO 1: Get an instance of the proxy
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := getSimpleRoutingQueueProxy(sdkConfig)

	log.Printf("Reading simple queue %s", d.Id())
	return gcloud.WithRetriesForRead(ctx, d, func() *resource.RetryError {
		/*
			CREATE-TODO 2: Call the proxy function getRoutingQueue(ctx context.Context, id string) to find our queue, passing in the ID from the resource data object
			The returned value are: Queue, Status Code (int), error
			If the error is not nil, we should pass the status code to the function gcloud.IsStatus404ByInt(int)
			If the status code is 404, return a resource.RetryableError. Otherwise, it should be a NonRetryableError
		*/
		currentQueue, respCode, err := proxy.getRoutingQueue(ctx, d.Id())
		if err != nil {
			if gcloud.IsStatus404ByInt(respCode) {
				return resource.RetryableError(fmt.Errorf("failed to read queue %s: %v", d.Id(), err))
			}
			return resource.NonRetryableError(fmt.Errorf("failed to read queue %s: %v", d.Id(), err))
		}

		// Define consistency checker
		cc := consistency_checker.NewConsistencyCheck(ctx, d, meta, ResourceSimpleRoutingQueue())

		// CREATE-TODO 3: Set our values in the schema resource data, based on the values in the Queue object returned from the API
		_ = d.Set("name", *currentQueue.Name)
		resourcedata.SetNillableValue(d, "calling_party_name", currentQueue.CallingPartyName)
		resourcedata.SetNillableValue(d, "enable_transcription", currentQueue.EnableTranscription)

		return cc.CheckState()
	})
}

// updateSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to update a simple queue in Genesys cloud.
func updateSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// CREATE-TODO 1: Get an instance of the proxy
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := getSimpleRoutingQueueProxy(sdkConfig)

	log.Printf("Updating simple queue %s", d.Id())

	// CREATE-TODO 2: Create variables for each field in our schema.ResourceData object
	name := d.Get("name").(string)
	callingPartyName := d.Get("calling_party_name").(string)
	enableTranscription := d.Get("enable_transcription").(bool)

	// CREATE-TODO	3: Create a queue struct using the Genesys Cloud platform go sdk
	queueUpdate := &platformclientv2.Queuerequest{
		Name:                &name,
		CallingPartyName:    &callingPartyName,
		EnableTranscription: &enableTranscription,
	}

	// CREATE-TODO 4: Call the proxy function updateRoutingQueue(context.Context, id string, *platformclientv2.Queuerequest) to update our queue
	_, _, err := proxy.updateRoutingQueue(ctx, d.Id(), queueUpdate)
	if err != nil {
		return diag.Errorf("failed to update queue %s: %v", name, err)
	}

	return readSimpleRoutingQueue(ctx, d, meta)
}

// deleteSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to delete a simple queue from Genesys cloud.
func deleteSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// CREATE-TODO 1: Get an instance of the proxy (done)
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := getSimpleRoutingQueueProxy(sdkConfig)

	log.Printf("Deleting simple queue %s", d.Id())

	/*
		CREATE-TODO 2: Call the proxy function deleteRoutingQueue(ctx context.Context, id string)
		Again, we won't be needing the returned response object
	*/
	if _, err := proxy.deleteRoutingQueue(ctx, d.Id()); err != nil {
		return diag.Errorf("failed to delete queue %s: %v", d.Id(), err)
	}

	// Check that queue has been deleted by trying to get it from the API
	time.Sleep(5 * time.Second)
	return gcloud.WithRetries(ctx, 30*time.Second, func() *resource.RetryError {
		_, respCode, err := proxy.getRoutingQueue(ctx, d.Id())

		if err == nil {
			return resource.NonRetryableError(fmt.Errorf("error deleting routing queue %s: %s", d.Id(), err))
		}
		if gcloud.IsStatus404ByInt(respCode) {
			// Success: Routing Queue deleted
			log.Printf("Deleted routing queue %s", d.Id())
			return nil
		}

		return resource.RetryableError(fmt.Errorf("routing queue %s still exists", d.Id()))
	})
}
