resource "podinate_pod" "my_pod" {
  account    = "...my_account..."
  id         = "hello-world"
  image      = "wordpress"
  name       = "Hello World"
  project_id = "hello-world"
  tag        = "6.0"
}