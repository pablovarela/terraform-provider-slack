---
page_title: "Provider: Slack"
---

# Slack Provider

The Slack provider is used to interact with Slack resources supported by Slack.
The provider needs to be configurred with a valid token before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
provider slack {
  token = var.slack_token
}

resource slack_conversation test {
  name       = "my-channel"
  topic      = "The topic for my channel"
  members    = []
  is_private = true
}
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the AWS
 `provider` block:

* `token` - (Mandatory) The Slack token. It must be provided,
but it can also be sourced from the `SLACK_TOKEN` environment variable.
