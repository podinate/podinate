# Shared Volumes
A Podinate shared volume is a shared persistent storage space that multiple pods can access at the same time. 

## Spec
This is a minimal demo of a shared volume in action. The `writer` pod will continually write to the volume, while the `reader` pod will read it out. To see the output, run `podinate -p shared-volume logs -f reader`. 
```hcl title="shared_volume.pcl"
project "shared-volume" {
    name = "Shared Volume"
    account_id = "default"
}
pod "reader" {
    name = "Shared Volume Reader"
    image = "ubuntu"
    command = [ "tail", "-f", "/data/time" ]
    project_id = "shared-volume"
    shared_volume {
        volume_id = "poddy"
        path = "/data"
    }
}
pod "writer" {
    name = "Shared Volume Writer"
    image = "ubuntu"
    command = [ "/bin/bash", "-c", "--" ]
    arguments = [ "while true; do date >> /data/time; sleep 2; done;" ]
    project_id = "shared-volume"
    shared_volume {
        volume_id = "poddy"
        path = "/data"
    }
}
shared_volume "poddy" {
    project_id = "shared-volume"
    size = 1
}
```