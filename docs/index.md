---
page_title: "Provider: Slack"
---

# Slack Provider

The Slack provider is used to interact with Slack resources supported by Slack.
The provider needs to be configurred with a valid token before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    slack = {
      source  = "pablovarela/slack"
      version = "~> 0.1"
    }
  }
  required_version = ">= 0.13"
}

# Configure Slack Provider
provider slack {
  token = var.slack_token
}

# Create a Slack channel
resource slack_conversation test {
  name       = "my-channel"
  topic      = "The topic for my channel"
  is_private = true
}
```

## Authentication

The Slack provider requires an Slack API token. It can be provided by different
means:

- Static token
- Environment variables

### Static Token

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system. 

A static can be provided by adding `token` in-line in the Slack provider block:

Usage:
```hcl
provider slack {
  token = var.slack_token
}
```
 
### Environment Variables
 
You can provide your token via the `SLACK_TOKEN` environment variable:

```hcl
provider slack {}
```

Usage:

```sh
$ export SLACK_TOKEN="my-token"
$ terraform plan
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the AWS
 `provider` block:

* `token` - (Mandatory) The Slack token. It must be provided,
but it can also be sourced from the `SLACK_TOKEN` environment variable.
