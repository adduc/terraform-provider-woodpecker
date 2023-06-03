package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataOrgSecret(t *testing.T) {
	name := "data.woodpecker_organization_secret.test_org_secret"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: orgSecretDataConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "owner", "test_org"),
					resource.TestCheckResourceAttr(name, "name", "test_org_secret"),
					resource.TestCheckNoResourceAttr(name, "value"),
					resource.TestCheckResourceAttr(name, "events.#", "1"),
					resource.TestCheckResourceAttr(name, "events.0", "push"),
				),
			},
		},
	})
}

const orgSecretDataConfig = `
resource "woodpecker_organization_secret" "test_org_secret" {
	owner = "test_org"
	name   = "test_org_secret"
	value  = "test_value"
	events = ["push"]
}

data "woodpecker_organization_secret" "test_org_secret" {
	owner = woodpecker_organization_secret.test_org_secret.owner
	name = woodpecker_organization_secret.test_org_secret.name
	depends_on = [woodpecker_organization_secret.test_org_secret]
}
`
