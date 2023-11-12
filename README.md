# podinate
A fast and easy way to get a project from development to prod

## Repo Structure
We use a monorepo, each top level folder here represents a single service, or holds some shared libraries

## Getting started with development

### Local Kubernetes environment
To get a complete working Podinate server, we need to get a local Kubernetes environment going:
#### Ubuntu:
```
sudo apt-get install k3d
```
#### Arch:
```
sudo pacman -S k3d
```

### Create a new Kubernetes cluster
Create a new Kubernetes cluster to run the code in
```
k3d cluster create podinate-dev
```
If you already have a k3d cluster running for something else, try this: 
```
k3d cluster create podinate-dev --api-port 6444
```
Check the cluster all looks okay. If the single node's status is "Ready" then so are you. 
```
$ kubectl get node 
NAME                        STATUS   ROLES                  AGE   VERSION
k3d-podinate-dev-server-0   Ready    control-plane,master   86d   v1.27.4+k3s1
```
You'll probably be typing `kubectl` a lot, so this is highly recommended.
```
$ alias k=kubectl
$ k get node
NAME                        STATUS   ROLES                  AGE   VERSION
k3d-podinate-dev-server-0   Ready    control-plane,master   86d   v1.27.4+k3s1
```


### Spin up API server and Postgres
First we spin up Postgres in our new kubernetes cluster: 
```
k create namespace api
k apply -f kubernetes/masterdb-postgres.yaml
```

### Load the SQL file into Postgres


### Install the API backend
