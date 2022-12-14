---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "woodpecker_organization_secret Resource - terraform-provider-woodpecker"
subcategory: ""
description: |-
  Provides a organization secret. For more
          information see Woodpecker CI's documentation https://woodpecker-ci.org/docs/usage/secrets
---

# woodpecker_organization_secret (Resource)

Provides a organization secret. For more 
		information see [Woodpecker CI's documentation](https://woodpecker-ci.org/docs/usage/secrets)

## Example Usage

```terraform
resource "woodpecker_organization_secret" "secret" {
  owner = "example_org"
  name  = "example secret"
  value = "example value"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Secret Name
- `owner` (String) Organization name
- `value` (String, Sensitive) Secret Value

### Optional

- `events` (Set of String) One or more event types where secret is available (one of push, tag, pull_request, deployment, cron, manual)
- `images` (Set of String) List of images where this secret is available, leave empty to allow all images
- `plugins_only` (Boolean) Whether secret is only available for plugins

### Read-Only

- `id` (Number) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# Syntax: <owner>/<name>
terraform import woodpecker_organization_secret.secret "example_org/test secret"
```
