---
subcategory: "Slack"
page_title: "Slack: slack_conversation"
---

# slack_conversation Data Source

Use this data source to get information about a Slack conversation for use in other
resources.

## Example Usage

```hcl
data slack_conversation test {
  channel_id = "my-channel"
}
```

## Argument Reference

The following arguments are supported:

* `channel_id` - (Required) The ID of the channel

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - name of the public or private channel.
* `members` - user IDs of users in the channel.
* `topic` - topic for the channel.
* `purpose` - purpose of the channel.
* `creator` - is the user ID of the member that created this channel.
* `created` - is a unix timestamp.
* `is_private` - means the conversation is privileged between two or more members.
* `is_archived` - indicates a conversation is archived. Frozen in time.
* `is_shared` - means the conversation is in some way shared between multiple workspaces.
* `is_ext_shared` - represents this conversation as being part of a Shared Channel
with a remote organization.
* `is_org_shared` - explains whether this shared channel is shared between Enterprise
Grid workspaces within the same organization.
* `is_general` - will be true if this channel is the "general" channel that includes
all regular team members.
