package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRepository_basic(t *testing.T) {
	name := "woodpecker_repository.test_repo"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: repositoryConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "owner", "test_user"),
					resource.TestCheckResourceAttr(name, "name", "test_repo"),
				),
			},
			// Import testing
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateId:     "test_user/test_repo",
				ImportStateVerify: true,
			},
			// Update/Read testing
			{
				Config: repositoryConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "owner", "test_user"),
					resource.TestCheckResourceAttr(name, "name", "test_repo"),
				),
			},
		},
	})
}

var repositoryConfig = `
resource "woodpecker_repository" "test_repo" {
	owner = "test_user"
	name = "test_repo"
}
`
