# Kubernetes Ingress with Podinate

The recommended way to do ingress with Podinate is using the [Nginx Ingress Controller](https://kubernetes.github.io/ingress-nginx/). Note that there are several different implementations of an "Nginx Ingress Controller". Podinate recommends the one published by the Kubernetes team unless a business wants to get official support from Nginx for the Ingress Controller published by F5. 

## Installation 
Follow the [official installation instructions](https://kubernetes.github.io/ingress-nginx/deploy/) to get started. They have a variety of installation guides for various platforms and clouds. 

## On-Prem / Homelab
In an on-prem / homelab environment, one option is to set up the Nginx with a NodePort service. This will be set up by default when installing the Nginx Ingress Controller. The details of the Nodeports can be found by describing the nginx Service object in Kubernetes.
```bash
kubectl -n ingress-nginx -n ingress-nginx describe service ingress-nginx-controller
```
For a homelab, forward the NodePorts indicated for the http and https ports from your router. If there's a need to get more fancy than that, have a look at [MetalLB](https://metallb.io/) for allocating IPs outside of the nodes. 

## Cloud
The popular clouds have their own sections in the [Nginx Ingress Controller Installation Guide](https://kubernetes.github.io/ingress-nginx/deploy/). If the cloud isn't mentioned, try doing a basic installation. Most cloud providers are configured so a *Service* of type LoadBalancer will get an actual external load balancer configured automatically, so this installation should work. 


## See Also
- [Nginx Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
- [Nginx Ingress Controller Installation](https://kubernetes.github.io/ingress-nginx/deploy/)
- [MetalLB](https://metallb.io/)