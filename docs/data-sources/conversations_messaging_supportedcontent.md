---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "genesyscloud_conversations_messaging_supportedcontent Data Source - terraform-provider-genesyscloud"
subcategory: ""
description: |-
  Genesys Cloud supported content data source. Select an supported content by name
---

# genesyscloud_conversations_messaging_supportedcontent (Data Source)

Genesys Cloud supported content data source. Select an supported content by name

## Example Usage

```terraform
data "genesyscloud_conversations_messaging_supportedcontent" "supported_content" {
  name = "Test Supported Content"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) supported content name

### Read-Only

- `id` (String) The ID of this resource.