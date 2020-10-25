---
subcategory: "Slack"
page_title: "Slack: slack_conversation"
---

# slack_conversation Resource

Manages a Slack channel

## Example Usage

```hcl
resource slack_conversation test {
  name              = "my-channel"
  topic             = "The topic for my channel"
  permanent_members = []
  is_private        = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) name of the public or private channel.
* `topic` - (Optional) topic for the channel.
* `purpose` - (Optional) purpose of the channel.
* `permanent_members` - (Optional) user IDs to add to the channel.
* `is_private` - (Optional) create a private channel instead of a public one.
* `is_archived` - (Optional) indicates a conversation is archived. Frozen in time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The channel ID (e.g. C015QDUB7ME).
* `creator` - is the user ID of the member that created this channel.
* `created` - is a unix timestamp.
* `is_shared` - means the conversation is in some way shared between multiple workspaces.
* `is_ext_shared` - represents this conversation as being part of a Shared Channel
with a remote organization.
* `is_org_shared` - explains whether this shared channel is shared between Enterprise
Grid workspaces within the same organization.
* `is_general` - will be true if this channel is the "general" channel that includes
all regular team members.
