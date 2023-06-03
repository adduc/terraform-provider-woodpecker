package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRepositorySecret_basic(t *testing.T) {
	name := "woodpecker_repository_secret.test_repo_secret"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: repositorySecretConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "repo_owner", "test_user"),
					resource.TestCheckResourceAttr(name, "repo_name", "test_repo"),
					resource.TestCheckResourceAttr(name, "name", "test_repo_secret"),
					resource.TestCheckResourceAttr(name, "events.#", "1"),
					resource.TestCheckResourceAttr(name, "events.0", "push"),
				),
			},
			// Import testing
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateId:           "test_user/test_repo/test_repo_secret",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
			// Update/Read testing
			{
				Config: repositorySecretConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "repo_owner", "test_user"),
					resource.TestCheckResourceAttr(name, "repo_name", "test_repo"),
					resource.TestCheckResourceAttr(name, "name", "test_repo_secret"),
					resource.TestCheckResourceAttr(name, "events.#", "1"),
					resource.TestCheckResourceAttr(name, "events.0", "push"),
				),
			},
		},
	})
}

var repositorySecretConfig = `
resource "woodpecker_repository" "test_repo" {
	owner = "test_user"
	name  = "test_repo"
}
resource "woodpecker_repository_secret" "test_repo_secret" {
	repo_owner = woodpecker_repository.test_repo.owner
	repo_name  = woodpecker_repository.test_repo.name
	name       = "test_repo_secret"
	value      = "test_value"
	events     = ["push"]
}
`
