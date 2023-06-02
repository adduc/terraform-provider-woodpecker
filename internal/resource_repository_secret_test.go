package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRepositorySecret_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: repositorySecretConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_repository_secret.test", "repo_owner", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_secret.test", "repo_name", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_secret.test", "name", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_secret.test", "events.#", "1"),
					resource.TestCheckResourceAttr("woodpecker_repository_secret.test", "events.0", "push"),
				),
			},
			// Import testing
			{
				ResourceName:            "woodpecker_repository_secret.test",
				ImportState:             true,
				ImportStateId:           "test/test/test",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
			// Update/Read testing
			{
				Config: repositorySecretConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_repository_secret.test", "repo_owner", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_secret.test", "repo_name", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_secret.test", "name", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_secret.test", "events.#", "1"),
					resource.TestCheckResourceAttr("woodpecker_repository_secret.test", "events.0", "push"),
				),
			},
		},
	})
}

var repositorySecretConfig = `
resource "woodpecker_repository" "test" {
	owner = "test"
	name  = "test"
}
resource "woodpecker_repository_secret" "test" {
	repo_owner = woodpecker_repository.test.owner
	repo_name  = woodpecker_repository.test.name
	name       = "test"
	value      = "test"
	events     = ["push"]
}
`
