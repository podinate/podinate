terraform {
  required_providers {
    podinate = {
      source  = "podinate/podinate"
      version = "0.0.1"
    }
  }
}

provider "podinate" {
  # Configuration options
    server_url = "http://localhost:3001/v0"
    api_key_auth = "puak-YINXecuzUzF82vFYEDsROZXH_KlrAW49gq2Rx_FtclNnbLcYvgdAnS-dLowoYkEVDHZg2Fi8rVSoWw1pF6JXoQ=="

}

resource "podinate_project" "wordpress_project" {
  account = "my-second-account"
  id      = "hello-world"
  name    = "Podinate Blog"
}

resource "podinate_pod" "my_pod" {
  account = podinate_project.wordpress_project.account
  id         = "hello-world"
  image      = "wordpress"
  name       = "Hello World"
  project_id = podinate_project.wordpress_project.id
  tag        = "6.0"
}