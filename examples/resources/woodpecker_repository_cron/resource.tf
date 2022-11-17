resource "woodpecker_repository" "repo" {
  owner = "example_user"
  name  = "woodpecker_test"
}

resource "woodpecker_repository_cron" "cron" {
  repo_owner = woodpecker_repository.repo.owner
  repo_name  = woodpecker_repository.repo.name
  name       = "terraform cron"
  schedule   = "@weekly"
}
