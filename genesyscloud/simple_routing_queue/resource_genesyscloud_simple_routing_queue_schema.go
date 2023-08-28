package simple_routing_queue

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gcloud "terraform-provider-genesyscloud/genesyscloud"
	registrar "terraform-provider-genesyscloud/genesyscloud/resource_register"
)

const resourceName = "genesyscloud_simple_routing_queue"

func SetRegistrar(l registrar.Registrar) {
	l.RegisterDataSource(resourceName, DataSourceSimpleRoutingQueue())
	l.RegisterResource(resourceName, ResourceSimpleRoutingQueue())
}

func ResourceSimpleRoutingQueue() *schema.Resource {
	return &schema.Resource{
		Description: "Genesys Cloud Simple Routing Queue",

		// TODO: Specify our our functions that we defined in resource_genesyscloud_simple_routing_queue.go for performing CRUD operations.
		// For example:
		// ReadContext: gcloud.ReadWithPooledClient(readSimpleRoutingQueue)
		CreateContext: gcloud.CreateWithPooledClient(createSimpleRoutingQueue),
		ReadContext:   gcloud.ReadWithPooledClient(readSimpleRoutingQueue),
		UpdateContext: gcloud.UpdateWithPooledClient(updateSimpleRoutingQueue),
		DeleteContext: gcloud.DeleteWithPooledClient(deleteSimpleRoutingQueue),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			/*
				TODO: Define the following three fields:
				1. "name"                 | type: string  | required | description: "The name of our routing queue."
				2. "calling_party_name"   | type: string  | optional | description: "The name to use for caller identification for outbound calls from this queue."
				3. "enable_transcription" | type: boolean | optional | description: "Indicates whether voice transcription is enabled for this queue."
				Note: The field "enable_transcription" is also Computed. This lets the provider know that the API will compute and return a value, should
					the user not specify one in the resource config.
			*/
			"name": {
				Description: "The name for our routing queue.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"calling_party_name": {
				Description: "The name to use for caller identification for outbound calls from this queue.\n",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enable_transcription": {
				Description: "Indicates whether voice transcription is enabled for this queue.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func DataSourceSimpleRoutingQueue() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for Genesys Cloud Simple Routing Queues.",
		// TODO: As above, specify the function dataSourceSimpleRoutingQueueRead as the ReadContext of this Resource object
		ReadContext: gcloud.ReadWithPooledClient(dataSourceSimpleRoutingQueueRead),

		Schema: map[string]*schema.Schema{
			/*
				TODO: Define the only field in our data source:
				"name" | type: string | required | description: "The name of our routing queue."
			*/
			"name": {
				Description: "The queue name.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}
