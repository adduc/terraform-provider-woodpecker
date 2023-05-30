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
				Config: `
				resource "woodpecker_repository" "test" {
					owner = "test"
					name = "test"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_repository.test", "owner", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository.test", "name", "test"),
				),
			},
		},
	})
}
