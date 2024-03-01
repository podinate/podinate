resource "podinate_pod" "my_pod" {
  account = "...my_account..."
  environment = [
    {
      key    = "DB_PASSWORD"
      secret = true
      value  = "abc1234"
    },
  ]
  id         = "hello-world"
  image      = "wordpress"
  name       = "Hello World"
  project_id = "hello-world"
  services = [
    {
      domain_name = "my-blog.podinate.app"
      name        = "my-blog"
      port        = 80
      protocol    = "http"
      target_port = 80
    },
  ]
  tag = "6.0"
  volumes = [
    {
      class      = "standard"
      mount_path = "/var/www/html"
      name       = "blog-data"
      size       = 10
    },
  ]
}