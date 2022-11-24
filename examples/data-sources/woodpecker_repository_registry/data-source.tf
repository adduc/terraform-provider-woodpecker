data "woodpecker_repository_registry" "registry" {
  repo_owner = "example_user"
  repo_name  = "woodpecker_test"
  address    = "docker.io"
}
