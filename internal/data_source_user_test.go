package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataUser(t *testing.T) {
	name := "data.woodpecker_user.test_user"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: userDataConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "login", "test_user"),
				),
			},
		},
	})
}

const userDataConfig = `
data "woodpecker_user" "test_user" {
	login = "test_user"
}
`
