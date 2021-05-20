---
subcategory: "Slack"
page_title: "Slack: slack_usergroup"
---

# slack_usergroup Data Source

Use this data source to get information about a usergroups for use in other
resources. The data source returns enabled groups only.

## Required scopes

This resource requires the following scopes:

- [usergroups:read](https://api.slack.com/scopes/usergroups:read)

The Slack API methods used by the resource are:

- [usergroups.list](https://api.slack.com/methods/usergroups.list)

If you get `missing_scope` errors while using this resource check the scopes against
the documentation for the methods above.

## Example Usage

```hcl
data slack_usergroup by_name {
  name = "my-usergroup"
}

data slack_usergroup by_id {
  usergroup_id = "USERGROUP00"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the usergroup
- `usergroup_id` - (Optional) The id of the usergroup

The data source expects exactly one of these fields, you can't set both.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the usergroup
- `description` - The short description of the User Group.
- `handle` - The mention handle.
- `users` - The user IDs that represent the entire list of users for the
  User Group.
- `channels` - The channel IDs for which the User Group uses as a default.
