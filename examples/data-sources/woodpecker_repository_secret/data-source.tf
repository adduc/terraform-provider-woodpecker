data "woodpecker_repository_secret" "secret" {
  repo_owner = "example_user"
  repo_name  = "woodpecker_test"
  name       = "example secret"
}
