---
page_title: "Provider: Woodpecker CI"
description: |-
  The Woodpecker CI provider provides utilities for automating
  configuration of Woodpecker CI instances.
---

# Woodpecker Provider

## Example Usage

```terraform
provider "woodpecker" {
  server = "https://woodpecker.example.com/"
  token  = "example-token"
}

## Resources

resource "woodpecker_repository" "repo" {
  owner  = "example_user"
  name   = "woodpecker_test"
  config = ".woodpecker.yml"
}

resource "woodpecker_repository_cron" "cron" {
  repo_owner = woodpecker_repository.repo.owner
  repo_name  = woodpecker_repository.repo.name
  name       = "terraform cron"
  schedule   = "@weekly"
}

## Data Sources

data "woodpecker_repository" "repo" {
  owner = woodpecker_repository.repo.owner
  name  = woodpecker_repository.repo.name
}

data "woodpecker_repository_cron" "cron" {
  repo_owner = woodpecker_repository_cron.cron.repo_owner
  repo_name  = woodpecker_repository_cron.cron.repo_name
  name       = woodpecker_repository_cron.cron.name
}

data "woodpecker_self" "self" {}

## Outputs

output "resource_repository" {
  value = woodpecker_repository.repo
}

output "resource_repository_cron" {
  value = woodpecker_repository_cron.cron
}

output "data_repository" {
  value = data.woodpecker_repository.repo
}

output "data_repository_cron" {
  value = data.woodpecker_repository_cron.cron
}

output "data_self" {
  value = data.woodpecker_self.self
}
```