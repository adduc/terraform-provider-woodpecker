data "woodpecker_repository_cron" "cron" {
  repo_owner = "example_user"
  repo_name  = "woodpecker_test"
  name       = "terraform cron"
}
