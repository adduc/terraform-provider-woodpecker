---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "woodpecker_repository_registry Data Source - terraform-provider-woodpecker"
subcategory: ""
description: |-
  Use this data source to get information on an existing registry for a repository
---

# woodpecker_repository_registry (Data Source)

Use this data source to get information on an existing registry for a repository

## Example Usage

```terraform
data "woodpecker_repository_registry" "registry" {
  repo_owner = "example_user"
  repo_name  = "woodpecker_test"
  address    = "docker.io"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `address` (String) Registry Address
- `repo_name` (String) Repository name
- `repo_owner` (String) User or organization responsible for repository

### Read-Only

- `email` (String) Registry Email
- `id` (Number) The ID of this resource.
- `token` (String, Sensitive) Registry Token
- `username` (String) Registry Username


