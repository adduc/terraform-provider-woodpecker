package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRepositoryRegistry_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: repositoryRegistryConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_repository_registry.test", "repo_owner", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_registry.test", "repo_name", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_registry.test", "address", "docker.io"),
					resource.TestCheckResourceAttr("woodpecker_repository_registry.test", "username", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_registry.test", "password", "test"),
				),
			},
			// Import testing
			{
				ResourceName:            "woodpecker_repository_registry.test",
				ImportState:             true,
				ImportStateId:           "test/test/docker.io",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update/Read testing
			{
				Config: repositoryRegistryConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_repository_registry.test", "repo_owner", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_registry.test", "repo_name", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_registry.test", "address", "docker.io"),
					resource.TestCheckResourceAttr("woodpecker_repository_registry.test", "username", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_registry.test", "password", "test"),
				),
			},
		},
	})
}

var repositoryRegistryConfig = `
resource "woodpecker_repository" "test" {
	owner = "test"
	name  = "test"
}
resource "woodpecker_repository_registry" "test" {
	repo_owner = woodpecker_repository.test.owner
	repo_name  = woodpecker_repository.test.name
	address    = "docker.io"
	username   = "test"
	password   = "test"
}
`
