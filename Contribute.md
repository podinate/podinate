# Contribute to Podinate
Thank you for your interest in contributing to Podinate. We are aiming to create an alternative container runtime and management engine that can be both a handy development environment on your laptop and a powerful production setup. 

If you're looking for something to contribute, check out our [Github Project](https://github.com/orgs/podinate/projects/1) for what we're currently up to. 


## Repo Structure
We use a monorepo, each top level folder here represents a single service, or holds some shared libraries. These are the most important ones:
- **api** - holds the OpenAPI definition for the Podinate API. 
- **api-backend**

 In the `api` folder is our API definition, if updated generate the updated client and server packages with the script `make api-generate`. It will ask you for a sudo password so it can update some permission issues from running the generator inside a Docker container. 

## Documentation
There is a Readme.md file inside important package folders. Please make sure to read `api-backend/Readme.md` and `api-backend/iam/Readme.md`.

## Getting started with development

### Local Kubernetes environment
To get a complete working Podinate server, we need to get a local Kubernetes environment going:
#### Ubuntu:
```
sudo apt-get install k3d
```
#### Arch:
```
sudo yay -S rancher-k3d-bin
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
You'll probably be typing `kubectl` a lot, so aliasing kubectl to just k is highly recommended.
```
$ alias k=kubectl
$ k get node
NAME                        STATUS   ROLES                  AGE   VERSION
k3d-podinate-dev-server-0   Ready    control-plane,master   86d   v1.27.4+k3s1
```


### Spin up API server and Postgres
First we spin up Postgres in our new kubernetes cluster: 
```
k create namespace podinate
k apply -f kubernetes/masterdb-postgres.yaml
```
Then install the API
```
k apply -f kubernetes/api-backend.yaml

```
This creates a pod running a hot-reload script for development. First upload the entire backend code, then I recommend using Kubycat to upload any changes into the pod during development. 
```
make dev-code-api
sudo npm install -g kubycat
kubycat ./kubycat.yaml # Leave this running to develop the backend 
```


To interact with it for development, forward port 3001 on your local machine to the API in the cluster
```
k -n podinate get pods
k -n podinate port-forward pods/api-backend-deployment-54c7b6895f-tg594 3001:3000 # Leave running to develop the backend
```


### Load the SQL file into Postgres
To get a postgres shell in the backend postgres instance, run the following
```
make postgres-shell
```
Then copy in the `database/masterdb.sql` file. 

## Use Insomnia
Load the `API/Insomnia.json` file into Insomnia to see the endpoints. First create an account, then a project, then create a pod inside the project. 

## Get Started Developing
Each top level folder has a Readme. Please read `api-backend/Readme.md` and `cli/Readme.md` to get started.  

The controller and cli communicate through client/server packages built off the OpenAPI spec in `api/`, if you make any changes to it, run `make api-generate` to rebuild it. 