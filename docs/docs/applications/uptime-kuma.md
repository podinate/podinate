# Uptime Kuma
[Uptime Kuma](https://uptime.kuma.pet/) is a self hosted monitoring service. Using Uptime Kuma, monitors for just about any service running on the internet can be configured, and alerts can be configured to notify the administrator when issues are detected. It provides several different types of monitoring, such as HTTP(S), ping, TCP Port check, DNS and Docker Host monitoring. It also supports many methods of notifying about issues it detects. 

Check out the [Run Your First App (Uptime Kuma)](../../getting-started/first-app) guide to install Uptime Kuma. 

## PodFile
The following PodFile runs a basic Uptime Kuma instance.
```hcl
podinate {
    package = "uptime-kuma"
    namespace = "uptime-kuma"
}

pod "uptime-kuma" {
    image = "louislam/uptime-kuma:1"
    service "kuma" {
        port = 80
        target_port = 3001
    }

    environment "PORT" {
        value = "3001"
    }
    
    volume "uptime-kuma" {
        mount_path = "/app/data"
        size = 1
    }
}
```

## Ingress
Ingress may not be needed for Uptime Kuma. Most likely it will only need to be accessed infrequently when you actually want to change some of the monitors. For this purpose, `kubectl port-forward` can be used like so:
```bash
kubectl -n uptime-kuma port-forward services/kuma 8080:
```
This will forward port 8080 on the local machine to Uptime Kuma whenever you need it, and when you don't need it the interface will be unavailable to the world. 

### CloudFlare Tunnels
If you're an existing CloudFlare customer, you may want to use [CloudFlare Tunnel](cloudflare-tunnel) to get access to services inside the cluster. 

## Monitoring Kubernetes Services
Uptime Kuma can monitor most Kubernetes services by accessing them at `http://service.namespace`. 

## Volumes
Uptime Kuma expects just one volume to be mounted at `/app/data` which is where the configuration and monitoring database will be stored. 

## See Also
- [Uptime Kuma Homepage](https://uptime.kuma.pet/)
- [Uptime Kuma Github](https://github.com/louislam/uptime-kuma)