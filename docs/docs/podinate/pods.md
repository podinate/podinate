# Podinate Pod
The Podinate Pod is the most commonly used resource provided by Podinate. It is a single container running in the cluster. 

## Create a Pod
This is the most basic Pod. It runs an nginx web server which will listen for requests.
```hcl title="pod.pcl"
pod "web-server" {
    image = "nginx"
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

## Shred Volumes
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

### Volumes
Pods can have one or more `volume` blocks giving them their own distinct storage block. 

```hcl
pod "volume-pod" {
    project_id = "ubuntu"
    image = "ubuntu"
    volume "storage" {
        
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