package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: userConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_user.test2", "login", "test2"),
				),
			},
			// Import testing
			{
				ResourceName:      "woodpecker_user.test2",
				ImportState:       true,
				ImportStateId:     "test2",
				ImportStateVerify: true,
			},
			// Update/Read testing
			{
				Config: userConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_user.test2", "login", "test2"),
				),
			},
		},
	})
}

var userConfig = `
resource "woodpecker_user" "test2" {
	login  = "test2"
}
  
`
