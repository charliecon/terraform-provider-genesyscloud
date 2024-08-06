package simple_routing_queue

import (
	"fmt"
	"strconv"
	"terraform-provider-genesyscloud/genesyscloud/provider"
	"terraform-provider-genesyscloud/genesyscloud/util"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSimpleRoutingQueue(t *testing.T) {
	var (
		resourceId          = "queue"
		name                = "Create 2023 Queue " + uuid.NewString()
		callingPartyName    = "Example Inc."
		enableTranscription = "true"

		fullResourcePath = fmt.Sprintf("%s.%s", resourceName, resourceId)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { util.TestAccPreCheck(t) },
		ProviderFactories: provider.GetProviderFactories(providerResources, providerDataSources),
		Steps: []resource.TestStep{
			{
				Config: generateSimpleRoutingQueueResource(
					resourceId,
					name,
					strconv.Quote(callingPartyName),
					enableTranscription,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourcePath, "name", name),
					resource.TestCheckResourceAttr(fullResourcePath, "calling_party_name", callingPartyName),
					resource.TestCheckResourceAttr(fullResourcePath, "enable_transcription", util.TrueValue),
				),
			},
		},
	})
}

func generateSimpleRoutingQueueResource(resourceId, name, callingPartyName, enableTranscription string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	name                 = "%s"
    calling_party_name   = %s
	enable_transcription = %s
}
`, resourceName, resourceId, name, callingPartyName, enableTranscription)
}
