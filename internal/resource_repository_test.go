package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRepository_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: repositoryConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_repository.test", "owner", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository.test", "name", "test"),
				),
			},
			// Import testing
			{
				ResourceName:      "woodpecker_repository.test",
				ImportState:       true,
				ImportStateId:     "test/test",
				ImportStateVerify: true,
			},
			// Update/Read testing
			{
				Config: repositoryConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_repository.test", "owner", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository.test", "name", "test"),
				),
			},
		},
	})
}

var repositoryConfig = `
resource "woodpecker_repository" "test" {
	owner = "test"
	name = "test"
}
`
