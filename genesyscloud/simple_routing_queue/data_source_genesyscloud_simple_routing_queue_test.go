package simple_routing_queue

import (
	"fmt"
	"terraform-provider-genesyscloud/genesyscloud/provider"
	"terraform-provider-genesyscloud/genesyscloud/util"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSimpleRoutingQueue(t *testing.T) {
	var (
		resourceId      = "queue_resource"
		dataSourceId    = "queue_data"
		simpleQueueName = "Create2023 queue " + uuid.NewString()

		fullPathToResource   = fmt.Sprintf("%s.%s", resourceName, resourceId)
		fullPathToDataSource = fmt.Sprintf("data.%s.%s", resourceName, dataSourceId)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { util.TestAccPreCheck(t) },
		ProviderFactories: provider.GetProviderFactories(providerResources, providerDataSources),
		Steps: []resource.TestStep{
			{
				Config: generateSimpleRoutingQueueResource(
					resourceId,
					simpleQueueName,
					util.NullValue,
					util.NullValue,
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
data "%s" "%s" {
	name = "%s"
	depends_on = [%s]
}
`, resourceName, dataSourceId, queueName, dependsOn)
}
