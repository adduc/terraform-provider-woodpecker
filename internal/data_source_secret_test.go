package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSecret(t *testing.T) {
	name := "data.woodpecker_secret.test_secret"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: secretDataConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "test_secret"),
					resource.TestCheckResourceAttr(name, "events.#", "1"),
					resource.TestCheckResourceAttr(name, "events.0", "push"),
				),
			},
		},
	})
}

const secretDataConfig = `
resource "woodpecker_secret" "test_secret" {
	name  = "test_secret"
	value = "test_value"
	events = ["push"]
}

data "woodpecker_secret" "test_secret" {
	name = "test_secret"
	depends_on = [woodpecker_secret.test_secret]
}
`
