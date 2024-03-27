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

Podinate uses Kubernetes to provide cluster services. If you just want to create a simple single-node cluster, the installer will do that for you. If you want to do something more custom, start with the [K3s Quickstart](https://docs.k3s.io/quick-start) documentation. If you want to install a cluster on multiple hosts we highly recommend the great [K3sup project](https://docs.k3s.io/quick-start).

### Install Podinate Cluster
The Podinate installer is designed to run on a fresh, dedicated Ubuntu 24.04 instance. This could be a virtual machine on your homelab, or a VM instance from your favourite cloud provider. The instance should have at least 2GB of ram. In the command prompt of your server instance, run:
```bash
curl -sfL https://raw.githubusercontent.com/podinate/podinate/main/kubernetes/install.sh | sudo bash
```
If Podinate detects an existing cluster, it will ask if you want to install Podinate cluster to that Kubernetes cluster.

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
This will show a table of your projects with only one entry: `Quick Start`. Let's set that as the default for your project.

```bash
echo "project: quick-start" > podinate.yaml
```
Now while you are in this directory, Podinate will look at the Project with the ID `quick-start` by default.

## Introducing Pods
A Podinate Pod is a container running in your cluster. You may be familiar with the concept of a Pod from other container managers like Kubernetes and Podman. 
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

