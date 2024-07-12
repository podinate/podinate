# Introduction to Kubernetes

Kubernetes is a container orchestration system for automating the management of containers, storage and configuration of applications. Kubernetes manages one or more computers into a cluster which can run container-based workloads. Kubernetes is one of the most widely deployed pieces of software in the world, meaning it is available on any cloud and can be deployed to any computer. 

## Concepts
Kubernetes provides many basic building blocks, but this introduction will focus on those most relevant to running containers and applications. 
- **Pods** - A Pod is a container running in the Kubernetes cluster. 
- **Control Plane** - The Kubernetes control plane provides an API which can be used to manage the rest of the system. 
- **Nodes** - A Node is a computer that is a member of the cluster, which the Control Plane can use to deploy workloads. 
- **Namespaces** - Nearly everything in Kubernetes is separated into Namespaces, which are a way to divide your cluster resources into logical units. For example, if you run a blog and a landing page on your cluster, you might have those two workloads separated into `blog` and `landing` namespaces. 
- **Services** - A Service provides an endpoint to connect to a group of Pods that provide one overall service. 
- **Volumes** - A volume is a persistent storage location to be attached to Pods. 

## Creating a Kubernetes Cluster
If this is your first time using Kubernetes or Podinate, try going through the [Podinate Quick Start](../getting-started/quick-start.md) guide to install a Kubernetes cluster and Podctl. 

## kubectl
Kubernetes provides an official CLI tool for interacting with the Kubernetes Control Plane, called Kubectl. Kubectl provides a handful of commands that can be applied to any type of *Object* in the Kubernetes cluster:
- **get [type]** - Show a list of all objects of the given type, for example `kubectl get pods`.
- **describe [type] [name]** - Get details of the given object, for example `kubectl describe pod my-app`. 
- **apply -f [file]** - Apply a file of YAML or JSON definitions of Kubernetes objects. They will be created or updated to match the desired state. 

### Kubectl and Namespaces
Most resources in Kubernetes are divided into Namespaces, which allows you to group related resources on the cluster. There is also a default Namespace, called `default`. 

Access objects in a specific namespace: 
```bash
kubectl -n my-app get pods
```
Keeping `-n my-app` at the start of the command allows you to easily edit the end of the command to change what kubectl is doing. 

Access objects in all namespaces: 
```bash
kubectl get -A pods
```
Unfortunately, -A needs to be placed after the sub-command name. 

You can also set which Namespace Kubectl will use if none is specified like so: 
```bash
kubectl config set-context --current --namespace=my-app
```
This is useful because most of the time you'll be interacting with one namespace for a long time, for example while deploying and debugging a particular app. 

## See Also
- [k3d](https://k3d.io/)