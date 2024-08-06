package simple_routing_queue

import (
	"terraform-provider-genesyscloud/genesyscloud/provider"
	resourceExporter "terraform-provider-genesyscloud/genesyscloud/resource_exporter"
	registrar "terraform-provider-genesyscloud/genesyscloud/resource_register"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const resourceName = "genesyscloud_simple_routing_queue"

func SetRegistrar(l registrar.Registrar) {
	l.RegisterResource(resourceName, ResourceSimpleRoutingQueue())
	l.RegisterDataSource(resourceName, DataSourceSimpleRoutingQueue())
}

func ResourceSimpleRoutingQueue() *schema.Resource {
	return &schema.Resource{
		Description: "Genesys Cloud Simple Routing Queue",

		// CREATE-TODO 1: Specify our our functions that we defined in resource_genesyscloud_simple_routing_queue.go for performing CRUD operations.
		// For example:
		// CreateContext: provider.CreateWithPooledClient(createSimpleRoutingQueue),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		// Here are the docs for the Schema struct: https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#Schema
		Schema: map[string]*schema.Schema{
			/*
				CREATE-TODO 2: Define the following three fields:
				1. "name"                 | Type: schema.TypeString  | Required | Description: "The name of the simple routing queue."
				2. "calling_party_name"   | Type: schema.TypeString  | Optional | Description: "The name to use for caller identification for outbound calls from this queue."
				3. "enable_transcription" | Type: schema.TypeBool    | Optional | Description: "Indicates whether voice transcription is enabled for this queue."
				Note: The field "enable_transcription" is also Computed. This lets the provider know that the API will compute and return a value, should
					the user not specify one in the resource config.

				An example field:
				"foo_bar": {
				    Description: "The foo for the bar.",
					Required:    true,
					Type:        schema.TypeString,
				},
			*/
		},
	}
}

func DataSourceSimpleRoutingQueue() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for Genesys Cloud Simple Routing Queues.",
		// CREATE-TODO 3: As above, specify the function dataSourceSimpleRoutingQueueRead as the ReadContext of this Resource object

		Schema: map[string]*schema.Schema{
			/*
				CREATE-TODO 4: Define the only field in our data source:
				"name" | Type: schema.TypeString  | Required | Description: "The name of the simple routing queue."
			*/
		},
	}
}

func RoutingQueueExporter() *resourceExporter.ResourceExporter {
	return &resourceExporter.ResourceExporter{
		GetResourcesFunc: provider.GetAllWithPooledClient(getAllRoutingQueues),
	}
}
