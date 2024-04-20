# Cert Manager on Kubernetes
[Cert manager](https://cert-manager.io/) is software to manage certificates and SSL on a Kubernetes or Podinate cluster. It issues certificates using the [ACME protocol](https://letsencrypt.org/how-it-works/) for completely automated certificate management and issuance, so there's no need to manually rotate certificates once per year like you may be used to. 

Cert-manager is installed by default when you set up a Podinate cluster. 

## Single-domain Certificates
To create an HTTPS website on Podinate, you'll need to set up Port Forwarding to forward Podinate's load balancer ports to the Podinate cluster. You can then create an HTTPS website using the following PCL: 
```hcl 
pod "https-test" {
    image = "nginx" 
    service "web-secure" {
        port = 80 #The port on the Pod
        protocol = "https"
        domain_name = "testsecure.example.com"
    }
}
```
It's a good idea to test that the domain points to your Podinate cluster before trying to issue an SSL certificate, or Let's Encrypt may rate limit your IP address. 
```hcl
pod "http-test" {
    image = "nginx" 
    service "web-insecure" {
        port = 80 #The port on the Pod
        protocol = "http"
        domain_name = "testsecure.example.com"
    }
}
```

## Wildcard Certificates 
If the above example of issuing a single-domain certificate is used, you may notice a lot of traffic coming to your new test website. This is because the certificate is published to Certificate Transparency Logs, which is scanned by many automated crawlers (not all of them bad). If you're planning to run many things on the Podinate cluster, it may be worthwhile to set up a Wildcard certificate for `*.example.com` instead of issuing them one by one. 


## See Also
- [Cert Manager Official Website](https://cert-manager.io/)
- [Cert Manager Documentation](https://cert-manager.io/docs/)
- [Let's Encrypt](https://letsencrypt.org/)