package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceOrganizationSecret_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: organizationSecretConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_organization_secret.test", "owner", "testorg"),
					resource.TestCheckResourceAttr("woodpecker_organization_secret.test", "name", "test"),
					resource.TestCheckResourceAttr("woodpecker_organization_secret.test", "events.#", "1"),
					resource.TestCheckResourceAttr("woodpecker_organization_secret.test", "events.0", "push"),
				),
			},
			// Import testing
			{
				ResourceName:            "woodpecker_organization_secret.test",
				ImportState:             true,
				ImportStateId:           "testorg/test",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
			// Update/Read testing
			{
				Config: organizationSecretConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_organization_secret.test", "owner", "testorg"),
					resource.TestCheckResourceAttr("woodpecker_organization_secret.test", "name", "test"),
					resource.TestCheckResourceAttr("woodpecker_organization_secret.test", "events.#", "1"),
					resource.TestCheckResourceAttr("woodpecker_organization_secret.test", "events.0", "push"),
				),
			},
		},
	})
}

var organizationSecretConfig = `
resource "woodpecker_organization_secret" "test" {
	owner = "testorg"
	name   = "test"
	value  = "test"
	events = ["push"]
}
`
