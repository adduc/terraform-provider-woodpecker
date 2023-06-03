package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRepositoryCron_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: repositoryCronConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_repository_cron.test", "repo_owner", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_cron.test", "repo_name", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_cron.test", "name", "test_cron"),
					resource.TestCheckResourceAttr("woodpecker_repository_cron.test", "schedule", "@daily"),
				),
			},
			// Import testing
			{
				ResourceName:      "woodpecker_repository_cron.test",
				ImportState:       true,
				ImportStateId:     "test/test/test_cron",
				ImportStateVerify: true,
			},
			// Update/Read testing
			{
				Config: repositoryCronConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("woodpecker_repository_cron.test", "repo_owner", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_cron.test", "repo_name", "test"),
					resource.TestCheckResourceAttr("woodpecker_repository_cron.test", "name", "test_cron"),
					resource.TestCheckResourceAttr("woodpecker_repository_cron.test", "schedule", "@daily"),
				),
			},
		},
	})
}

var repositoryCronConfig = `
resource "woodpecker_repository" "test" {
	owner = "test"
	name  = "test"
}
resource "woodpecker_repository_cron" "test" {
	repo_owner = woodpecker_repository.test.owner
	repo_name  = woodpecker_repository.test.name
	name    = "test_cron"
	schedule = "@daily"
}
`
