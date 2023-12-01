
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
	api_key_auth = var.podinate_api_key
}

resource "podinate_project" "wordpress_project" {
	account = var.account_id
	id      = replace(lower(var.project_name), " ", "-")
	name    = var.project_name
}

resource "podinate_pod" "wordpress_pod" {
	account = podinate_project.wordpress_project.account
	project_id = podinate_project.wordpress_project.id
	id         = "wordpress"
	image      = "wordpress"
	name       = "WordPress"
	tag        = "6"
}

resource "podinate_pod" "database_pod" {
	account = podinate_project.wordpress_project.account
	project_id = podinate_project.wordpress_project.id
	id         = "mariadb"
	image      = "mariadb"
	name       = "MySQL"
	tag        = "11"
}
		
