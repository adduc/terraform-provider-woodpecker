package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceOrganizationSecret_basic(t *testing.T) {
	name := "woodpecker_organization_secret.test_org_secret"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: organizationSecretConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "owner", "test_org"),
					resource.TestCheckResourceAttr(name, "name", "test_org_secret"),
					resource.TestCheckResourceAttr(name, "events.#", "1"),
					resource.TestCheckResourceAttr(name, "events.0", "push"),
				),
			},
			// Import testing
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateId:           "test_org/test_org_secret",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
			// Update/Read testing
			{
				Config: organizationSecretConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "owner", "test_org"),
					resource.TestCheckResourceAttr(name, "name", "test_org_secret"),
					resource.TestCheckResourceAttr(name, "events.#", "1"),
					resource.TestCheckResourceAttr(name, "events.0", "push"),
				),
			},
		},
	})
}

var organizationSecretConfig = `
resource "woodpecker_organization_secret" "test_org_secret" {
	owner = "test_org"
	name   = "test_org_secret"
	value  = "test_value"
	events = ["push"]
}
`
