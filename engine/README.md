# Engine
The engine is the beating heart of Podinate. It can take a .pf (hopefully soon a .yaml) file and make the Kubernetes state match that state. 

## Podfile
A Podfile is an HCL file with a bunch of resources declared in it. Those resources are: 
- Pod
    - Service
    - Ingress
    - Volume
    - Shared Volume 
- Shared Volume
- Podinate 
    - Contains information about the package and the namespace