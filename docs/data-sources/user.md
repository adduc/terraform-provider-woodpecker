---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "woodpecker_user Data Source - terraform-provider-woodpecker"
subcategory: ""
description: |-
  Use this data source to get information on an existing user
---

# woodpecker_user (Data Source)

Use this data source to get information on an existing user

## Example Usage

```terraform
data "woodpecker_self" "self" {
  login = "username"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `login` (String) Username for user

### Read-Only

- `active` (Boolean) Whether user is active in the system
- `admin` (Boolean) Whether user is a Woodpecker admin
- `avatar` (String) Avatar URL for user
- `email` (String) Email address for user
- `id` (Number) User ID


