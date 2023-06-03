package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRepositoryRegistry_basic(t *testing.T) {
	name := "woodpecker_repository_registry.test_repo_registry"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: repositoryRegistryConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "repo_owner", "test_user"),
					resource.TestCheckResourceAttr(name, "repo_name", "test_repo"),
					resource.TestCheckResourceAttr(name, "address", "docker.io"),
					resource.TestCheckResourceAttr(name, "username", "reg_test_user"),
					resource.TestCheckResourceAttr(name, "password", "reg_test_pass"),
				),
			},
			// Import testing
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateId:           "test_user/test_repo/docker.io",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update/Read testing
			{
				Config: repositoryRegistryConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "repo_owner", "test_user"),
					resource.TestCheckResourceAttr(name, "repo_name", "test_repo"),
					resource.TestCheckResourceAttr(name, "address", "docker.io"),
					resource.TestCheckResourceAttr(name, "username", "reg_test_user"),
					resource.TestCheckResourceAttr(name, "password", "reg_test_pass"),
				),
			},
		},
	})
}

var repositoryRegistryConfig = `
resource "woodpecker_repository" "test_repo" {
	owner = "test_user"
	name  = "test_repo"
}
resource "woodpecker_repository_registry" "test_repo_registry" {
	repo_owner = woodpecker_repository.test_repo.owner
	repo_name  = woodpecker_repository.test_repo.name
	address    = "docker.io"
	username   = "reg_test_user"
	password   = "reg_test_pass"
}
`
