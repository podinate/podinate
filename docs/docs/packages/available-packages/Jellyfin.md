# Jellyfin on Kubernetes / Podinate 
[Jellyfin](https://jellyfin.org/) is an open source media center that started as a community-focused fork of Emby. It can be used to show a collection of media for example, images, videos and Linux ISOs. 

## Install
The following PCL will give a basic Jellyfin setup: 
```hcl title="jellyfin.pcl"
project "media" {
    name = "Media"
    account_id = "default"
}

pod "jellyfin" {
    project_id = "media"
    
    name = "Jellyfin"
    image = "jellyfin/jellyfin"
    
    volume "config" {
        size = 50
        mount_path = "/config"
        class = "nvme"
    }

    volume "cache" {
        size = 100
        mount_path = "/cache"
    }
}
```