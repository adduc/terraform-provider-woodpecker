resource "woodpecker_repository" "repo" {
  owner = "example_user"
  name  = "woodpecker_test"
}

resource "woodpecker_repository_registry" "registry" {
  repo_owner = woodpecker_repository.repo.owner
  repo_name  = woodpecker_repository.repo.name
  address    = "docker.io"
  username   = "exampleusername"
  password   = "examplepassword"
}
