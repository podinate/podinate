variable "podinate_api_key" {
    type = string 
    description = "API key for Podinate"
}

variable "account_id" {
    type = string
    description = "Account ID for Podinate"
}

variable "project_name" {
    type = string
    default = "WordPress Blog"
    description = "Project name for Podinate"
}