# Kubernetes Ingress

Kubernetes Ingress is used to expose services running on the cluster to the outside world. The most popular ways to provide ingress are to use the Nginx ingress controller or to use Traefik as the controller. Most cloud providers will provide their own solution, but these are not covered here. 

The recommended way to do ingress with Podinate is using the [Nginx Ingress Controller](https://kubernetes.github.io/ingress-nginx/). Note that there are several different implementations of an "Nginx Ingress Controller". Podinate recommends the one published by the Kubernetes team unless a business wants to get official support from Nginx for the Ingress Controller published by F5. 

## PodFile
Once Nginx is installed, all that is needed to add ingress to a Pod is to add a service and ingress block to the Pod's specification:
```hcl
service {
    port = 80
    target_port = 8080 // Optional
    ingress {
        hostname = "tunnel.example.com
    }
}
```
For example, the following PodFile will run a container hosting a clone of the 1024 game, and send any requests for `ingress.example.com` to the Pod:
```hcl 
podinate {
    package = "2048"
    namespace = "2048"
}

pod "game" {
    image = "alexwhen/docker-2048"
    service "game" {
        port = 80
        ingress {
            hostname = "ingress.example.com"
        }
    }
}
```

## Installation 
Follow the [official installation instructions](https://kubernetes.github.io/ingress-nginx/deploy/) to get started. They have a variety of installation guides for various platforms and clouds. For most other cases, the default installation will work. 

## Certificates
See [Certificates](certificates.md) to set up Let's Encrypt and certificate management for the Ingress Controller. Alternatively, if you're already a CloudFlare customer, consider using [CloudFlare Tunnels](../applications/cloudflare-tunnel.md) to provide access from outside the cluster and handle certificate issuance.

## On-Prem / Homelab
In an on-prem / homelab environment, one option is to set up the Nginx with a NodePort service. This will be set up by default when installing the Nginx Ingress Controller. The details of the Nodeports can be found by describing the nginx Service object in Kubernetes.
```bash
kubectl -n ingress-nginx -n ingress-nginx describe service ingress-nginx-controller
```
Forwarding port 80 and 443 on the router to the indicated ports will make the Nginx ingress controller available to the world. 

Another option for small environment is to use Cloudflare Tunnels to handle ingress. Install the Nginx ingress controller and then set up a [Cloudflare Tunnel pod](../applications/cloudflare-tunnel.md). Set the Cloudflare Tunnels to forward all traffic to `http://ingress-nginx-controller.ingress-nginx` and let Cloudflare handle SSL Certificates. 

## Cloud
The popular clouds have their own sections in the [Nginx Ingress Controller Installation Guide](https://kubernetes.github.io/ingress-nginx/deploy/). If the cloud isn't mentioned, try doing a basic installation. Most cloud providers are configured so a *Service* of type LoadBalancer will get an actual external load balancer configured automatically, so this installation should work. 

## See Also
- [Nginx Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
- [Nginx Ingress Controller Installation](https://kubernetes.github.io/ingress-nginx/deploy/)
- [MetalLB](https://metallb.io/)