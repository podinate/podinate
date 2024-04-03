# Pod 
Pod is the package that does the most interaction with the Kubernetes API. I've put most of the kubernetes logic for pods themselves into a separate file kubernetes.go.

## Glossary
### Pod 
Pods are Kubernetes deployments under the hood, allowing us to scale them up and down etc. 

### Environment Variables
Each pod can have environment variables set so that the user can configure various parts of their app. 

### Services
Services are a wrapper around the Kubernetes services object, with some ingress thrown in on top. Currently only if a hostname is set, the service will also have an ingress created. 

## Standards
All methods for interacting with Kubernetes should be on the Pod struct called `ensureX()`, for example `ensureNamespace()` and `ensureDeployment()`. 