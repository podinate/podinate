# Getting Started
This guide will take you through installing Podinate and a Kubernetes cluster on your local machine. You will then spin up Uptime Kuma and an Nginx Pod to learn how to Kubernet the Podinate way. 

## Install the CLI
Podinate is available through Homebrew for both Mac and Linux. If you don't have Homebrew, run the command on the [Homebrew homepage](https://brew.sh/) to install it. 
```bash
brew install podinate/tap/podinate k3d
```
This will install the Podinate CLI and K3d, which you will use to create a local Kubernetes cluster. 

You will also need Docker installed, if you don't have it already, check [Install Docker Desktop Mac](https://docs.docker.com/desktop/install/mac-install/)

## Create Local Cluster
Before anything can be run, a Kubernetes cluster needs to be created. For this tutorial, K3d will be used. 
```bash
k3d cluster create local
```
Now check the cluster has installed correctly:
```bash
kubectl get pods -A
```
The output should look like: 
```
NAMESPACE     NAME                                      READY   STATUS      RESTARTS   AGE
kube-system   coredns-6799fbcd5-9474v                   1/1     Running     0          26s
kube-system   local-path-provisioner-6c86858495-kkbdd   1/1     Running     0          26s
kube-system   metrics-server-54fd9b65b-9rqwd            0/1     Running     0          26s
kube-system   helm-install-traefik-crd-np6qg            0/1     Completed   0          26s
kube-system   helm-install-traefik-4lg85                0/1     Completed   1          26s
kube-system   svclb-traefik-f4e950dc-84xzw              2/2     Running     0          9s
kube-system   traefik-f4564c4f4-pgwgj                   0/1     Running     0          9s
```

## Run an Ubuntu Pod
First, let's create an Ubuntu Pod we can play with. First let's create a directory to hold this tutorial. 
```bash
mkdir podinate-quick-start
cd podinate-quick-start
```

Now copy the following into `ubuntu.pod`.
```hcl title="ubuntu.pod"
podinate {
    package = "ubuntu"
    namespace = "default"
}

pod "ubuntu" {
    image = "ubuntu"
    tag = "latest" 
    command = [ "/bin/bash", "-c", "--" ]
    arguments = [ "while true; do echo 'Hello from Podinate!'; sleep 2; done;" ]
}
```
This file creates two things. At the top, it creates a Project called Quick Start, then it creates a Pod called `ubuntu`, which runs the latest Ubuntu image, and runs a certain command in a loop. We'll learn more about Projects and Pods in the next sections. 

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
podinate logs -f ubuntu
```
You'll now be seeing "Hello from Podinate!" logged every two seconds. This is a very useful tool to see what is going on in the Pod. The `-f` means to keep following the log and show any new lines that come up. You can omit it if you just want to see the current contents of the log. 

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