# Engine
The engine is the beating heart of Podinate. It can take a .pf or Kubernetes manifest file and make the Kubernetes state match that state. 

## Podfile
A Podfile is an HCL file with a bunch of resources declared in it. Those resources are: 
- Pod
    - Service
        - Ingress
            - Annotations
    - Volume
    - Shared Volume 
    - Environment
    - Resources
- Shared Volume
- Podinate 
    - Contains information about the package and the namespace

## Packages
First, the engine scans one or more files and creates one or more *packages* from them. These packages contain an array of type Resource, and an array of Pods and SharedVolumes. They also contain a default Namespace. 

The Resource type is essentially a wrapper around an array of Kubernetes runtime.Objects, with a couple of other defaults like a name and a type to display to the user. 

I'm currently in the process of ripping out the Pods and SharedVolumes array, in favour of just returning arrays of resources when scanning the files. 

## Plan 
A plan is how Podinate needs to get from the current state to the desired state. Most of this is internal, as Kubernetes lets just just yeehaw any object into its API and it will try to get to that state. Have a look at package.Plan in plan.go. It basically creates the default namespace if set and then loops over the package.Resources slice. For each it runs GetResourceChangeForResource which is a big method that needs to be renamed to GetObjectChangeForObject once the Pod and SharedVolume specific changes are ripped out. 

### GetResourceChangeForObject 
This is the function that does all the Kubernetes interaction magic. First it gets a Magic Kubernetes RestHelper for the object. It first attempts to grab the current object from the Kubernetes API. If the object doesn't exist, it will do a Kubernetes Dry Run of creating the object. This will determine if there's any issue with the configuration before we attempt to create it. If there's any issue, it will return an error here and return to the user. This is in contrast to Kubectl which just treats objects one by one and attempts to yeehaw them into the cluster. 

If the object already exists, the method will grab some metadata about the object and transfer it to the new one. This makes sense for Kubernetes annotations, which are often added after the fact by various controllers. However, for labels, we should be adding proper support for them to Podinate. This means if the user updated some labels in Kube, Podinate won't remove them. This makes sense for a PodFile, where a user might add a label to make something in the cluster behave a certain way, and Podinate shouldn't remove that the next time it's updated. However, this doesn't make much sense for Kubernetes YAML files. That remains an unsolved problem. 

Finally, the function does a DeepCompare of the object. If it's the same, it will return a NoOp change. If not, it will dry run replacement of the object, which will allow the Kube to tell us if there's any issue with the new spec or any reason it can't be updated. 

## Helpers 
There's a sub-package called helpers which does some things that are commonly needed when interacting with Kubernetes. For example, there's functions to convert between the different Kubernetes generic object types that different packages use. 