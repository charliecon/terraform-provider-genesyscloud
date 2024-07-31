package simple_routing_queue

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"terraform-provider-genesyscloud/genesyscloud/provider"
	"terraform-provider-genesyscloud/genesyscloud/util"
	"testing"
)

func TestAccDataSourceSimpleRoutingQueue(t *testing.T) {
	var (
		resourceId      = "queue_resource"
		dataSourceId    = "queue_data"
		simpleQueueName = "Create2023 queue " + uuid.NewString()

		fullPathToResource   = resourceName + "." + resourceId
		fullPathToDataSource = "data." + resourceName + "." + dataSourceId
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { util.TestAccPreCheck(t) },
		ProviderFactories: provider.GetProviderFactories(nil, nil),
		Steps: []resource.TestStep{
			{
				Config: generateSimpleRoutingQueueResource(
					resourceId,
					simpleQueueName,
					"null",
					"null",
				) + generateSimpleRoutingQueueDataSource(
					dataSourceId,
					simpleQueueName,
					fullPathToResource,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						fullPathToDataSource, "id",
						fullPathToResource, "id",
					),
				),
			},
		},
	})
}

func generateSimpleRoutingQueueDataSource(dataSourceId, queueName, dependsOn string) string {
	return fmt.Sprintf(`
data "genesyscloud_simple_routing_queue" "%s" {
	name = "%s"
	depends_on = [%s]
}
`, dataSourceId, queueName, dependsOn)
}
