---
subcategory: "Slack"
page_title: "Slack: slack_user"
---

# change_inventory_item Data Source

Use this data source to get information about an inventory item for use in other
resources.

## Example Usage

```hcl
data slack_user test {
  name = "my-user"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the user
