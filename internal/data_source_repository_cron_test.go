package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataRepoCron(t *testing.T) {
	name := "data.woodpecker_repository_cron.test_repo_cron"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: NewProto6ProviderFactory(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: repoCronDataConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "repo_owner", "test_user"),
					resource.TestCheckResourceAttr(name, "repo_name", "test_repo"),
					resource.TestCheckResourceAttr(name, "name", "test_cron"),
					resource.TestCheckResourceAttr(name, "schedule", "@daily"),
				),
			},
		},
	})
}

const repoCronDataConfig = `
resource "woodpecker_repository" "test_repo" {
	owner = "test_user"
	name  = "test_repo"
}

resource "woodpecker_repository_cron" "test_repo_cron" {
	repo_owner = woodpecker_repository.test_repo.owner
	repo_name  = woodpecker_repository.test_repo.name
	name    = "test_cron"
	schedule = "@daily"
}

data "woodpecker_repository_cron" "test_repo_cron" {
	repo_owner = woodpecker_repository_cron.test_repo_cron.repo_owner
	repo_name  = woodpecker_repository_cron.test_repo_cron.repo_name
	name	   = woodpecker_repository_cron.test_repo_cron.name
}
`
