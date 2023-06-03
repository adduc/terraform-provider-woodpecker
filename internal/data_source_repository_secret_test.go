package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataRepoSecret(t *testing.T) {
	name := "data.woodpecker_repository_secret.test_secret"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: repoSecretDataConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "repo_owner", "test_user"),
					resource.TestCheckResourceAttr(name, "repo_name", "test_repo"),
					resource.TestCheckResourceAttr(name, "name", "test_secret"),
					resource.TestCheckResourceAttr(name, "events.#", "1"),
					resource.TestCheckResourceAttr(name, "events.0", "push"),
				),
			},
		},
	})
}

const repoSecretDataConfig = `
resource "woodpecker_repository" "test_repo" {
	owner = "test_user"
	name  = "test_repo"
}

resource "woodpecker_repository_secret" "test_secret" {
	repo_owner = woodpecker_repository.test_repo.owner
	repo_name  = woodpecker_repository.test_repo.name
	name       = "test_secret"
	value      = "test_value"
	events     = ["push"]
}

data "woodpecker_repository_secret" "test_secret" {
	repo_owner = woodpecker_repository.test_repo.owner
	repo_name  = woodpecker_repository.test_repo.name
	name       = woodpecker_repository_secret.test_secret.name
	depends_on = [woodpecker_repository_secret.test_secret]
}
`
