# Podinate Configuration Language

Podinate does away with the hundred-line long YAML files associated with other container orchestration solutions. It uses Podinate Configuration Language (PCL) which uses the Hashicorp Configuration Language syntax you may be familiar with from Terraform. With PCL, you can define various aspects of your application, such as containers, volumes, services, and dependencies, in a concise and intuitive manner.

## Features
### Simplified Syntax

PCL utilizes a clean and straightforward syntax based on HCL, allowing users to define application configurations with ease. The intuitive structure makes it easy to understand and modify configurations as needed.

### Declarative Definitions
With PCL, you can declare your application's desired state, specifying the resources and dependencies it requires for deployment. PCL abstracts away the complexities of Kubernetes configurations, providing a higher-level interface for defining applications.

### Modular and Reusable

PCL supports modularity and reusability, allowing you to define reusable components and modules that can be easily integrated into different applications. This promotes code reuse and maintainability across your projects.

### Integration with Podinate

PCL seamlessly integrates with Podinate, allowing you to directly apply your configuration files to deploy and manage your applications within the Podinate environment. This tight integration streamlines the deployment process and ensures compatibility with Podinate's features and capabilities.

## Example

Here's an example of a simple application defined in PCL: 

```hcl title="web-server.pcl"
project "file-server" {
    name = "File Server"
    account_id = "default"
}

# Nginx http server
pod "nginx" {
    name = "Files Archive"
    image = "nginx"
    tag = "latest" 
    shared_volume {
        source = "files"
        path = "/var/www/html"
    }

    # Expose the web server on port 30080 of the cluster nodes
    service "web" {
        type = "http"
        node_port = 30080
    }
}

shared_volume "files" {
    size = 10 
}

```
This PCL file will create three things: 

- A [Project](project) called `File Server`
- A [Pod](pod) called `Files Archive` - which has a shared volume and a service
- A [Shared Volume](shared-volume) with the ID `files` which is 10GB in size