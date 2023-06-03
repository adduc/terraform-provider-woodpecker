package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceUser_basic(t *testing.T) {
	name := "woodpecker_user.test_user_2"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: userConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "login", "test_user_2"),
				),
			},
			// Import testing
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateId:     "test_user_2",
				ImportStateVerify: true,
			},
			// Update/Read testing
			{
				Config: userConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "login", "test_user_2"),
				),
			},
		},
	})
}

var userConfig = `
resource "woodpecker_user" "test_user_2" {
	login  = "test_user_2"
}
`
