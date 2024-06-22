# Certificates (HTTPS) on Kubernetes

When exposing services to end users from a Kubernetes cluster, getting an SSL certificate for that service is almost mandatory. For most small - medium Kubernetes installations, the easiest way to get one is to use [Cert Manager](https://cert-manager.io/) with Let's Encrypt.

## Installation
Once the [Nginx Ingress Controller](ingress) is installed, install Cert Manager with Helm. See also the Cert Manager Helm installation guide. 
```bash
helm repo add jetstack https://charts.jetstack.io --force-update

helm install \
    cert-manager jetstack/cert-manager \
    --namespace cert-manager \
    --create-namespace \
    --version v1.15.0 \
    --set crds.enabled=true
```
### Create Let's Encrypt Issuer
To issue publicly trusted certificates, a *ClusterIssuer* needs to be created. The most popular way to do this is by using Let's Encrypt to issue certificates.
```yaml title="lets-encrypt-issuer.yaml"
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-issuer
  annotations:
    podinate.com/default-cluster-issuer: "true"
spec:
  acme:
    email: your_email_address
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-issuer-account-key
    solvers:
      - http01:
          ingress:
            ingressClassName: nginx
```

### Set the Default Issuer for Cert-Manager
By setting the default issuer, you don't need to specify how each Kubernetes Ingress should obtain a certificate, they will use Let's Encrypt by default.
```bash
helm upgrade \
    cert-manager jetstack/cert-manager \
    --namespace cert-manager \
    --create-namespace \
    --version v1.15.0 \
    --set crds.enabled=true --set ingressShim.defaultIssuerName=letsencrypt-issuer --set ingressShim.defaultIssuerKind=letsencrypt-issuer
```


## Debugging
Any issues with the cert-manager configuration will mostly show in the cert-manager logs. These logs can be shown as follows:
```bash
kubectl -n cert-manager logs -f -l app=cert-manager
```

## See Also
- 
- [Cert Manager Helm Installation](https://cert-manager.io/docs/installation/helm/)