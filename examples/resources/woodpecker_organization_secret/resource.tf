resource "woodpecker_organization_secret" "secret" {
  owner = "example_org"
  name  = "example secret"
  value = "example value"
}
