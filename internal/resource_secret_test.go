package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecret_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: secretConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_secret.test", "name", "test"),
					resource.TestCheckResourceAttr("woodpecker_secret.test", "value", "test"),
				),
			},
			// Import testing
			{
				ResourceName:            "woodpecker_secret.test",
				ImportState:             true,
				ImportStateId:           "test",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
			// Update/Read testing
			{
				Config: secretConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_secret.test", "name", "test"),
					resource.TestCheckResourceAttr("woodpecker_secret.test", "value", "test"),
				),
			},
		},
	})
}

var secretConfig = `
resource "woodpecker_secret" "test" {
	name  = "test"
	value = "test"
	events = ["push"]
}
  
`
