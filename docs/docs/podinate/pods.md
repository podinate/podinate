# Podinate Pod
The Podinate Pod is the most commonly used resource provided by Podinate. It is a single container running in the cluster. 

## Create a Pod
This is the most basic Pod. It runs an nginx web server which will listen for requests.
```hcl title="pod.pcl"
pod "web-server" {
    image = "nginx"
}
```

## Environment Variables
Environment variables can be set in the Pod by adding an `environment` block.
```hcl title="environment_variable.pcl"
pod "web-server" {
    image = "nginx"
    environment "SERVER_NAME" {
        value = "example.com"
    }
}
```

## Services
Pods that run services should have a `service` block to take advantage of Kubernetes service discovery and optional ingress. 
```hcl title="service.pcl"
pod "web-server" {
    image = "nginx"
    environment "SERVER_NAME" {
        value = "example.com"
    }

    service "www" {
        port = 80
        protocol = "http"

        // Adding protocol = http and a domain_name will set up ingress on this service
        domain_name = "example.com"
    }
}
```


## Volumes
A volume is a separate persistent storage space, dedicated to the Pod. 

Let's say we want this nginx to have a space to use for a persistent cache. 
```hcl title="volume.pcl"
pod "web-server" {
    image = "nginx"
    volume "cache" {
        path = "/cache"
    }
}
```

## Shared Volumes
A Shared Volume is also a persistent storage space, but it can be shared among many pods. An app can use it to store persistent data that may need to be used by multiple Pods. 
```hcl title="shared_volume.pcl"
shared_volume "files" { // Create shared volume
    size = 2
}
pod "web-server" {
    image = "nginx"
    volume "cache" {
        path = "/cache"
    }
    shared_volume { // Attach the shared volume to the Pod
        volume_id = "files"
        path = "/var/www/html"
    }
}
```

## Spec
This is the complete spec of a Pod Block:
```hcl title="pod_full.pcl"
pod "ubuntu" {
    project_id = "ubuntu"
    name = "My Ubuntu Pod"
    image = "ubuntu"
    tag = "22.04"
    command = [ "/bin/bash", "-c", "--" ]
    arguments = [ "while true; do echo 'Hello from Podinate!'; sleep 2; done;" ]
}
```

{{ read_csv('./pod_reference.csv') }}

### Volume Spec
Pods can have one or more `volume` blocks giving them their own distinct storage block. 

```hcl
pod "volume-pod" {
    project_id = "ubuntu"
    image = "ubuntu"
    volume "storage" {
        // The path in the Pod to mount the volume
        // Required
        path = "/var/www/html" 

        // Size limit of the Volume in GB 
        size = 20 
        
        // The Kubernetes storage class to use for the device
        // To see available storage classes run "podinate kubectl get storageclasses"
        class = "ssd"
    }
}
```

## Roadmap
For alpha, you can create a single long-runnning Pod running a single container. Podinate plans to expand this to allow you to run a group of identical containers soon.

Podinate plans to support three types of container:
```hcl
pod "service" {
    // type = "service" // Default, optional
}
pod "job" {
    // A pod that runs once and exits once done
    type = "job"
}
pod "cronjob" {
    // A pod that is run at a set interval
    type = "cron"
    schedule = "*/5 * * *"
}
```

## Kubernetes
A Podinate Pod maps to a Kubernetes Stateful Set 

## See Also