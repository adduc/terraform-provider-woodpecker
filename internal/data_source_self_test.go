package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSelf(t *testing.T) {
	name := "data.woodpecker_self.test_self"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: selfDataConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "login", "test_user"),
				),
			},
		},
	})
}

const selfDataConfig = `
data "woodpecker_self" "test_self" {}
`
