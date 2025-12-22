---
subcategory: "Slack"
page_title: "Slack: slack_usergroups"
---

# slack_usergroups Data Source

Use this data source to get information about all usergroups. This is useful for
filtering usergroups in Terraform rather than making multiple API calls. The data
source returns enabled groups only.

## Required scopes

This resource requires the following scopes:

- [usergroups:read](https://api.slack.com/scopes/usergroups:read)

The Slack API methods used by the resource are:

- [usergroups.list](https://api.slack.com/methods/usergroups.list)

If you get `missing_scope` errors while using this resource check the scopes against
the documentation for the methods above.

## Example Usage

```hcl
data "slack_usergroups" "all" {}

# Filter usergroups by handle prefix
locals {
  engineering_groups = [
    for g in data.slack_usergroups.all.usergroups : g
    if startswith(g.handle, "eng-")
  ]
}

# Use with for_each
resource "slack_conversation" "team_channels" {
  for_each = { for g in data.slack_usergroups.all.usergroups : g.id => g }
  name     = "team-${each.value.handle}"
}
```

## Argument Reference

This data source has no arguments.

## Attribute Reference

The following attributes are exported:

- `usergroups` - A list of all usergroups. Each usergroup has the following attributes:
  - `id` - The ID of the usergroup
  - `name` - The name of the usergroup
  - `description` - The short description of the User Group
  - `handle` - The mention handle
  - `users` - The user IDs that represent the entire list of users for the User Group
  - `channels` - The channel IDs for which the User Group uses as a default
