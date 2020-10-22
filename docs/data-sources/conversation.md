---
subcategory: "Slack"
page_title: "Slack: slack_conversation"
---

# change_inventory_item Data Source

Use this data source to get information about an inventory item for use in other
resources.

## Example Usage

```hcl
data slack_conversation test {
  name = "my-channel"
}
```

## Argument Reference

The following arguments are supported:

* `channel_id` - (Required) The ID of the channel

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - Name of the public or private channel.
* `topic` - Topic for the channel.
* `members` -  Members to add to the channel.
* `is_private` - Whether the channel is private or not.
