package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecret_basic(t *testing.T) {
	name := "woodpecker_secret.test_secret"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: secretConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "test_secret"),
					resource.TestCheckResourceAttr(name, "value", "test_value"),
				),
			},
			// Import testing
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateId:           "test_secret",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
			// Update/Read testing
			{
				Config: secretConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "test_secret"),
					resource.TestCheckResourceAttr(name, "value", "test_value"),
				),
			},
		},
	})
}

var secretConfig = `
resource "woodpecker_secret" "test_secret" {
	name  = "test_secret"
	value = "test_value"
	events = ["push"]
}
`
