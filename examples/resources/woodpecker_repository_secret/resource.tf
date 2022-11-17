resource "woodpecker_repository" "repo" {
  owner  = "example_user"
  name   = "woodpecker_test"
}

resource "woodpecker_repository_secret" "secret" {
  repo_owner = woodpecker_repository.repo.owner
  repo_name  = woodpecker_repository.repo.name
  name       = "example secret"
  value      = "example value"
}