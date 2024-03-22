# Getting Started

This guide aims to get you started with Podinate and start familiarising yourself with the platform. We will install the Podinate CLI on your local machine, set up a server, and run an example workload. 

## Install the CLI
Podinate CLI is what you will use to interact with the Podinate server. It is currently available through Homebrew for both Mac and Linux. Let's get it installed.
```bash
brew install Podinate/tap/podinate
```
This will install Podinate CLI from our Homebrew tap. 

Verify that the CLI was successfully installed:
```bash
podinate version
```
The output will look like:
```
Podinate CLI vX.X.X
```

## Create a Cluster


If you're creating a new cluster from scratch, we recommend a standard Ubuntu 22.04 installation as a starting point. If you don't have an existing Kubernetes cluster you want to use, run the following command to install the recommended K3s cluster. 
```bash
curl -sfL https://get.k3s.io | sh -
```
For more information on K3s visit the [K3s Quickstart](https://docs.k3s.io/quick-start) documentation. If you want to install a cluster on multiple hosts we highly recommend the great [K3sup project](https://docs.k3s.io/quick-start).

### Install Podinate controller
Now we will install Podinate onto the newly created cluster. In the command prompt of your server instance, run:
```bash
curl -Lo https://raw.github.com/Podinate/podinate/.../server.sh | sudo bash
```
Podinate will use the standard $KUBECONFIG environment variable to connect to the cluster. If you want to use a specific kubeconfig file and context within the file, use the following command.
```bash
curl -Lo https://raw.github.com/Podinate/podinate/.../server.sh | KUBECONFIG=~/.kube/config CONTEXT=<your-context> sudo bash
```

??? note
    This will install:

    - Certbot manager 
    - Podinate Postgres database
    - Set up Postgres tables
    - Podinate controller
    - (future) ask about storage
    - Create default Podinate account 
    - Create admin Podinate user
    - Install the admin credentials to the local machine

## Let's Podinate!
<!-- You can now use Podinate as you might Docker. Most commands are the same. For example `podinate build` will run a build in the Podinate cluster and cache the file locally.  -->

First, let's create an Ubuntu Pod we can play with. First let's create a directory to hold this tutorial. 
```bash
mkdir podinate-quick-start
cd podinate-quick-start
```

Now copy the following into `ubuntu.pod`.
```hcl title="ubuntu.pod"
project "quick-start" {
    name = "Quick Start"
}

pod "ubuntu" {
    name = "Quick Start Ubuntu"
    image = "ubuntu"
    project_id = "quick-start"
}
```
We can now create our Ubuntu Pod by running;
```bash
podinate install ubuntu.pod
```
The process should only take a second, and now we have a running Ubuntu Pod. 

## Introducing Projects
Podinate divides your pods and other resources into Projects. This means you can keep resources organised logically by what they are a part of. Let's take a look at all our Projects now. 
```bash
podinate get projects
```
This will show a pretty table of your projects with only one entry: `Quick Start`. Let's set that as the default for your project.

```bash
echo "project: quick-start" > podinate.yaml
```
Now while you are in this directory, Podinate will look at the Project with the ID `quick-start` by default.

## Introducing Pods
A Podinate Pod is a container running in your cluster. You may be familiar with the concept of a Pod from other container managers like Kubernetes and Podman, and the concept is similar here. 
```bash
podinate get pods 
```
This will show a table of your pods, you should see only one called `Quick Start Ubuntu`, running the ubuntu:latest image. 

## Get Ubuntu Shell
We can now get a shell on the ubuntu pod by running the following command; 
```bash
podinate shell ubuntu
```
The `podinate shell` command is a convenient way to get inside of a container and debug. We can now run commands like we would on any Ubuntu system:
```bash
echo "Hello"
ping podinate.com -c 5
curl https://api64.ipify.org
```

