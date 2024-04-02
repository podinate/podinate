# Getting Started

This guide aims to get you started with Podinate and start familiarising yourself with the platform. We will install the Podinate CLI on your local machine, set up a server, and run an example workload. 


## Create a Cluster

Podinate uses Kubernetes to provide cluster services. If you just want to create a simple single-node cluster, the installer will do that for you. You can then add more nodes later. If you want to do something more custom, start with the [K3s Quickstart](https://docs.k3s.io/quick-start) documentation. If you want to install a cluster on multiple hosts we highly recommend the great [K3sup project](https://docs.k3s.io/quick-start).

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

### Install the CLI (Optional)
The Podinate server installer will set up the Podinate command line client for the root user on the server. If you want to be able to control your Podinate server from your local command line, the CLI is available through Homebrew for both Mac and Linux. If you don't have Homebrew, run the command on the [Homebrew homepage](https://brew.sh/) to install it. 
```bash
brew tap podinate/tap
brew install podinate
```
This will install Podinate CLI from our Homebrew tap. 

### Login to Podinate
The server installer will set up the Podinate credentials for the root user, and print out the credentials file at the end of the installation process. If you want to use Podinate as another user, or from your local machine, you can add the server by running: 
```bash
podinate login
```
Paste the credentials file, then press `control + s` to save the new profile.

### Run an Ubuntu Pod

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
podinate apply ubuntu.pod
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
This will show a table of your pods, you should see only one called `Quick Start Ubuntu`, running the `ubuntu:latest` image. 

### Check Pod Logs 
The Pod logs will contain the output of the program running inside the container. In this case, we didn't specify one so the default Entrypoint from the Dockerfile is used. 
```bash
podinate logs ubuntu
```
This is a very useful tool to see what is going on in the Pod. 

### Run Command in Pod
Podinate can run any command inside of our pod. This command will let us list the contents of the `/var` directory, for example. 
```bash
podinate exec ubuntu -- ls /var
```

### (Coming Soon) Get Ubuntu Shell
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