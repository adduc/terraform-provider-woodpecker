terraform {
  required_providers {
    woodpecker = {
      source  = "terraform.local/adduc/woodpecker"
      version = "0.0.1-dev"
    }
  }
}

provider "woodpecker" {
  server = "http://ci.172.17.0.1.nip.io/"
  token  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXh0IjoiamxvbmciLCJ0eXBlIjoidXNlciJ9.In1kQ3Idy57r-JPRjMSwslkVTFtMuflfe4zhIRX39Ws"
}

##
## Data Source: Repository
##

# data "woodpecker_repository" "repository" {
#   owner = "jlong"
#   name  = "repo-3"
# }

# output "repository" {
#   value = data.woodpecker_repository.repository
# }


##
## Data Source: Self
##

# data "woodpecker_self" "self" {}
# output "self" {
#   value = data.woodpecker_self.self
# }


##
## Resource: Repository
##

resource "woodpecker_repository" "repository" {
  owner = "jlong"
  name  = "repo-3"
  config = "b"
}

output "repository" {
  value = woodpecker_repository.repository
}