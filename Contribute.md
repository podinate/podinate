# Contribute to Podinate
Thank you for your interest in contributing to Podinate. Podinate is aiming to create an alternative container runtime and management engine that can be both a handy development environment on your laptop and a powerful production setup. 

If you're looking for something to contribute, check out our [Github Project](https://github.com/orgs/podinate/projects/1) and see what interests you.


## Repo Structure
The Podinate code uses a monorepo, each top level folder here represents a single service, or holds some shared libraries. These are the most important ones:
- **api** - holds the OpenAPI definition for the Podinate API. 
- **controller** - this is the backend that implements the Podinate API. It is installed into a Kubernetes cluster using the Kubernetes YAML in `kubernetes/controller.yaml`.
- **database** - the database schema for the controller. 
- **docs** - the mkdocs site at [docs.podinate.com](https://docs.podinate.com), probably the most important folder.
- **kubernetes** - Kubernetes YAML for installing Podinate and supporting services.
- **lib** - shared libraries, currently just the generated OpenAPI client for the OpenAPI spec. 

 In the `api` folder is our API definition, if updated generate the updated client and server packages with the script `make api-generate`. It will ask you for a sudo password so it can update some permission issues from running the generator inside a Docker container. 

## Documentation
There is a Readme.md file inside important package folders. Please make sure to read `controller/Readme.md` and `controller/iam/Readme.md`.

## Getting started with development

### Creating Development Environment
To get a complete working Podinate server, you need to get a local Kubernetes environment. For this Podinate recommends using K3d which you can install like so:
#### Arch:
```bash
sudo yay -S rancher-k3d-bin
```

#### curl | sudo bash
```bash
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
```

You will also need the Node tool Kubycat, which will sync the controller code into the Kubernetes cluster so you can easily develop it. 
```bash
npm install -g kubycat

## May be needed on some systems
apt-get install fswatch 
```

### Clone the Git Repository
Of course, you have to check out the git repository first. 
```bash
# Using ssh key
git clone git@github.com:podinate/podinate.git
# Using https
git clone https://github.com/podinate/podinate.git
```

Then `cd` into the new directory:
```bash
cd podinate
```

### Create a New Development Cluster
Once K3d is installed, there is a Make script to create a development cluster and deploy the Podinate controller to it. Run that now. 
```bash
make dev-cluster
```
The script will run through various steps, and will pause at the database migrations step to confirm you want to apply it. Press enter for apply and let the process continue. 

This script will create a folder called `testapp` and put the credentials file in there. 

From here I recommend creating three new tabs in your terminal / IDE. Run the following commands in three of them, and the fourth will be your work terminal to run `podinate` commands. 

In the first terminal tab, run:
```bash
make dev-code-upload
```
This will upload all of the code in `controller` to the running controller pod, and start Kubycat (installed earlier) to sync up any code changes as you make them. 

In the second terminal, run:
```bash
make dev-port-forward
```
This will open a port forward from your local machine to the controller inside the Kubernetes cluster. 

In the third terminal, run:
```bash
make dev-backend-logs 
```
This will follow the logs of what's happening in the backend Pod. For the most part you won't need to look at the other two tabs during development, but will probably refer to this one a lot. 

Finally, in the last terminal tab, run the following:
```bash
cd testapp
alias podinate="go run ../cli"
alias p=podinate
cat credentials.yaml | podinate login
podinate get projects 
```

You can now develop Podinate! 

## Get Started Developing
Each top level folder has a Readme. Please read `controller/Readme.md` and `cli/Readme.md` to get started.

If you're looking for something to start with, check out our [Github Project](https://github.com/orgs/podinate/projects/1) and see what interests you. 

The controller and cli communicate through client/server packages built off the OpenAPI spec in `api/`, if you make any changes to it, run `make api-generate` to rebuild it. 