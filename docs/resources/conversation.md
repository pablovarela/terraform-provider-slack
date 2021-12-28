---
subcategory: "Slack"
page_title: "Slack: slack_conversation"
---

# slack_conversation Resource

Manages a Slack channel

## Required scopes

This resource requires the following scopes:

- [channels:read](https://api.slack.com/scopes/channels:read) (public channels)
- [channels:write](https://api.slack.com/scopes/channels:write) (public channels)
- [groups:read](https://api.slack.com/scopes/groups:read) (private channels)
- [groups:write](https://api.slack.com/scopes/groups:write) (private channels)

The Slack API methods used by the resource are:

- [conversations.create](https://api.slack.com/methods/conversations.create)
- [conversations.setTopic](https://api.slack.com/methods/conversations.setTopic)
- [conversations.setPurpose](https://api.slack.com/methods/conversations.setPurpose)
- [conversations.info](https://api.slack.com/methods/conversations.info)
- [conversations.members](https://api.slack.com/methods/conversations.members)
- [conversations.kick](https://api.slack.com/methods/conversations.kick)
- [conversations.invite](https://api.slack.com/methods/conversations.invite)
- [conversations.rename](https://api.slack.com/methods/conversations.rename)
- [conversations.archive](https://api.slack.com/methods/conversations.archive)
- [conversations.unarchive](https://api.slack.com/methods/conversations.unarchive)

If you get `missing_scope` errors while using this resource check the scopes against
the documentation for the methods above.

## Example Usage

```hcl
resource "slack_conversation" "test" {
  name              = "my-channel"
  topic             = "The topic for my channel"
  permanent_members = []
  is_private        = true
}
```

```hcl
resource "slack_conversation" "nonadmin" {
  name              = "my-channel01"
  topic             = "The channel won't be archived on destroy"
  permanent_members = []
  is_private        = true
  action_on_destroy = "none"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) name of the public or private channel.
- `topic` - (Optional) topic for the channel.
- `purpose` - (Optional) purpose of the channel.
- `permanent_members` - (Optional) user IDs to add to the channel.
- `is_private` - (Optional) create a private channel instead of a public one.
- `is_archived` - (Optional) indicates a conversation is archived. Frozen in time.
- `action_on_destroy` - (Optional, Default `archive`) indicates whether the conversation should be archived or left
behind on destroy. Valid values are `archive | none`. Note that when set to `none` the conversation will be left as it
is  and as a result any subsequent runs of terraform apply with the same name  will fail.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The channel ID (e.g. C015QDUB7ME).
- `creator` - is the user ID of the member that created this channel.
- `created` - is a unix timestamp.
- `is_shared` - means the conversation is in some way shared between multiple workspaces.
- `is_ext_shared` - represents this conversation as being part of a Shared Channel
with a remote organization.
- `is_org_shared` - explains whether this shared channel is shared between Enterprise
Grid workspaces within the same organization.
- `is_general` - will be true if this channel is the "general" channel that includes
all regular team members.

## Import

`slack_conversation` can be imported using the ID of the conversation/channel, e.g.

```shell
terraform import slack_conversation.my_conversation C023X7QTFHQ
```
