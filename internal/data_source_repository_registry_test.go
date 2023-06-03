package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataRepoRegistry(t *testing.T) {
	name := "data.woodpecker_repository_registry.test_repo_registry"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: repoRegistryDataConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "repo_owner", "test_user"),
					resource.TestCheckResourceAttr(name, "repo_name", "test_repo"),
					resource.TestCheckResourceAttr(name, "address", "docker.io"),
					resource.TestCheckResourceAttr(name, "username", "reg_test_user"),
					resource.TestCheckNoResourceAttr(name, "password"),
				),
			},
		},
	})
}

const repoRegistryDataConfig = `
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

data "woodpecker_repository_registry" "test_repo_registry" {
	repo_owner = woodpecker_repository_registry.test_repo_registry.repo_owner
	repo_name  = woodpecker_repository_registry.test_repo_registry.repo_name
	address    = woodpecker_repository_registry.test_repo_registry.address
	depends_on = [woodpecker_repository_registry.test_repo_registry]
}
`
