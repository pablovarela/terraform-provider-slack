---
subcategory: "Slack"
page_title: "Slack: slack_user"
---

# slack_user Data Source

Use this data source to get information about a user for use in other
resources.

## Example Usage

```hcl
data slack_user by_name {
  name = "my-user"
}

data slack_user by_email {
  email = "my-user@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the user
* `email` - (Optional) The email of the user

The data source expects exactly one of these fields, you can't set both.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the user
