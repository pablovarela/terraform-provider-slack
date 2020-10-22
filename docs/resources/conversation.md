---
subcategory: "Slack"
page_title: "Slack: slack_conversation"
---

# slack_conversation Resource

Manages a Slack channel

## Example Usage

```hcl
resource slack_conversation test {
  name       = "my-channel"
  topic      = "The topic for my channel"
  members    = []
  is_private = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the public or private channel.
* `topic` - (Optional) Topic for the channel.
* `members` - (Optional) Members to add to the channel.
* `is_private` - (Optional) Create a private channel instead of a public one.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The channel ID (e.g. C015QDUB7ME).
