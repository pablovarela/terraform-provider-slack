---
subcategory: "Slack"
page_title: "Slack: slack_usergroup"
---

# slack_usergroup Resource

Manages a Slack User Group.

## Required scopes

This resource requires the following scopes:

- [usergroups:write](https://api.slack.com/scopes/usergroups:write)
- [usergroups:read](https://api.slack.com/scopes/usergroups:read)

The Slack API methods used by the resource are:

- [usergroups.create](https://api.slack.com/methods/usergroups.create)
- [usergroups.enable](https://api.slack.com/methods/usergroups.enable)
- [usergroups.disable](https://api.slack.com/methods/usergroups.disable)
- [usergroups.update](https://api.slack.com/methods/usergroups.update)
- [usergroups.list](https://api.slack.com/methods/usergroups.list)
- [usergroups.users.update](https://api.slack.com/methods/usergroups.users.update)

If you get `missing_scope` errors while using this resource check the scopes against
the documentation for the methods above.

## Example Usage

```hcl
resource "slack_usergroup" "my_group" {
  name        = "TestGroup"
  handle      = "test"
  description = "Test user group"
  users       = ["USER00"]
  channels    = ["CHANNEL00"]
}
```

Note that if a channel is removed from the `channels` list users are
**not** removed from the channel. In order to keep the users in the
groups and in the channel in sync set `permanent_users` in the channel:

```hcl
resource "slack_usergroup" "my_group" {
  name        = "TestGroup"
  handle      = "test"
  description = "Test user group"
  users       = ["USER00"]
}

resource "slack_conversation" "test" {
  name              = "my-channel"
  topic             = "The topic for my channel"
  permanent_members = slack_usergroup.my_group.users
  is_private        = true
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) a name for the User Group. Must be unique among User Groups.
- `description` - (Optional) a short description of the User Group.
- `handle` - (Optional) a mention handle. Must be unique among channels, users
  and User Groups.
- `users` - (Optional) user IDs that represent the entire list of users for the
  User Group.
- `channels` - (Optional) channel IDs for which the User Group uses as a default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The usergroup ID

## Import

`slack_usergroup` can be imported using the ID of the group, e.g.

```shell
terraform import slack_usergroup.my_group S022GE79E9G
```
