package simple_routing_queue

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	gcloud "terraform-provider-genesyscloud/genesyscloud"
	"time"
)

/*
   The data_source_genesyscloud_simple_routing_queue.go contains the data source implementation
   for the resource.
*/

// dataSourceSimpleRoutingQueueRead retrieves by search term the id in question
func dataSourceSimpleRoutingQueueRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// CREATE-TODO 1: . Get an instance of our proxy
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := getSimpleRoutingQueueProxy(sdkConfig)

	// CREATE-TODO 2: Grab our queue name from the schema.ResourceData object
	name := d.Get("name").(string)

	log.Printf("Finding queue with name '%s'", name)
	return gcloud.WithRetries(ctx, 15*time.Second, func() *resource.RetryError {
		// CREATE-TODO 3: Call to the proxy function getRoutingQueueIdByName(context.Context, string)
		// This function returns values in the following order: queueId (string), retryable (bool), err (error)
		queueId, retryable, err := proxy.getRoutingQueueIdByName(ctx, name)

		// CREATE-TODO 4: If the error is not nil, and retryable equals false, return a resource.NonRetryableError
		// to let the user know that an error occurred
		if err != nil && !retryable {
			return resource.NonRetryableError(fmt.Errorf("error finding queue '%s': %v", name, err))
		}

		// CREATE-TODO 5: If retryable equals true, return a resource.RetryableError and let them know the queue could not be found with that name
		if retryable {
			return resource.RetryableError(fmt.Errorf("no queue found with name '%s'", name))
		}

		// CREATE-TODO 6: If we made it this far, we can set the queue ID in the schema.ResourceData object, and return nil
		d.SetId(queueId)
		return nil
	})
}
