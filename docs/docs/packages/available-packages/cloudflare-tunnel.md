# Cloudflare Tunnel
[Cloudflare Tunnels](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/) are a way to connect your Podinate cluster to the internet via Cloudflare. The tunnel creates a connection out from a Pod in a Podinate cluster to Cloudflare's datacenters, so there's no need to configure port forwarding or certificates. 

## Install
Create a Cloudflare Tunnel by following [Cloudflare's Create Tunnel Documentation](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/get-started/create-remote-tunnel/). 

!!! important
    In step 2 of the Cloudflare documentation, specify `http://traefik.kube-system` to forward all requests into Podinate's load balancer.

You can use the following PCL to install the Cloudflare Tunnel connector on the Podinate cluster:
```hcl
project "cloudflare-tunnel" {
    name = "Cloudflare Tunnel"
    account_id = "default"
}

pod "tunnel" {
    name = "Cloudflare Tunnel"
    project_id = "cloudflare-tunnel"
    image = "cloudflare/cloudflared"
    arguments = [
        "tunnel",
        "--no-autoupdate",
        "run",
        "--token",
        "<your-token-here>"
    ]
    
}
```
Replace <your-token-here> with the long random string from the Cloudflare tunnel dashboard. 

## Send traffic to your Pod
In the Pod you want to expose through Cloudflare, add the following:
```hcl
service {
    type = "http"
    hostname = "www.example.com"
}
```

## See Also
- [Official Documentation](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/)
- [Create Tunnel - Cloudflare Documenation](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/get-started/create-remote-tunnel/)