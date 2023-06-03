package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataRepository(t *testing.T) {
	name := "data.woodpecker_repository.test_repo"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: repositoryDataConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "owner", "test_user"),
					resource.TestCheckResourceAttr(name, "name", "test_repo"),
				),
			},
		},
	})
}

const repositoryDataConfig = `
resource "woodpecker_repository" "test_repo" {
	owner = "test_user"
	name = "test_repo"
}

data "woodpecker_repository" "test_repo" {
	owner = "test_user"
	name = "test_repo"
	depends_on = [woodpecker_repository.test_repo]
}
`
