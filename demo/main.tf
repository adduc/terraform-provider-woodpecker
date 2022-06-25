terraform {
  required_providers {
    woodpecker = {
      source  = "jlong-ryzen-desktop/adduc/woodpecker"
      version = "0.0.1-dev"
    }
  }
}

provider "woodpecker" {
  server = "http://ci.172.17.0.1.nip.io/"
  token  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXh0IjoiamxvbmciLCJ0eXBlIjoidXNlciJ9.In1kQ3Idy57r-JPRjMSwslkVTFtMuflfe4zhIRX39Ws"
}

# data "woodpecker_repository" "repo" {
#   owner = "jlong"
#   name  = "2nd-repo"
# }

# output "repo" {
#   value = data.woodpecker_repository.repo
# }


data "woodpecker_self" "self" {}
output "self" {
  value = data.woodpecker_self.self
}
