# Getting Started
This guide will take you through installing Podinate and a Kubernetes cluster on your local machine. You will run a hello world container, and learn how to look at the logs and run commands to debug it.

## Install the CLI
Podinate is available for Windows through WSL, Mac through HomeBrew, or for Debian and Arch based Linux distros. 

=== "Windows"
    Podinate is available for Windows through Windows Subsystem for Linux. 
    
    If you haven't already, install WSL and Ubuntu by opening PowerShell and running the following command: 
    ```powershell
    wsl --install
    ```
    You will be asked to create a username and password for the new Linux subsystem, then it will automatically be started.

    With WSL installed, open the Ubuntu application and run the following commands to install Podinate. 
    ```bash
    wget --quiet -O - https://get.podinate.com/deb/gpg.key | sudo tee /etc/apt/keyrings/podinate.asc
    echo "deb [signed-by=/etc/apt/keyrings/podinate.asc] https://get.podinate.com/deb stable main" | sudo tee /etc/apt/sources.list.d/podinate.list
    sudo apt-get update
    sudo apt-get install podinate
    podinate version
    ```
=== "Mac (Homebrew)"
    Podinate is available through Homebrew for both Mac and Linux. If you don't have Homebrew, run the command on the [Homebrew homepage](https://brew.sh/) to install it. 
    ```bash
    brew install podinate/tap/podinate k3d
    ```
    This will install the Podinate CLI and K3d, which you will use to create a local Kubernetes cluster. You will also need Docker installed, if you don't have it already, check [Install Docker Desktop Mac](https://docs.docker.com/desktop/install/mac-install/)
=== "Linux (Debian)"
    Podinate can be installed from our Debian package repo. Podinate has no dependencies so should install on any Debian-based distro.
    Use the following commands to install Poinate cli:
    ```bash
    wget --quiet -O - https://get.podinate.com/deb/gpg.key | sudo tee /etc/apt/keyrings/podinate.asc
    echo "deb [signed-by=/etc/apt/keyrings/podinate.asc] https://get.podinate.com/deb stable main" | sudo tee /etc/apt/sources.list.d/podinate.list
    sudo apt-get update
    sudo apt-get install podinate
    podinate version
    ```
=== "Linux (Arch)"
    Podinate can be installed from the Arch User Repository using your favourite AUR helper.
    ```bash
    yay -S podinate
    podinate version
    ```

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
These are the Pods used to provide system services. It doesn't matter how many or what they do, just check they're all `Running` or `Completed`. This might take a few minutes. 

## Run a Hello World Pod
First, let's create an Ubuntu Pod we can play with. First let's create a directory to hold this tutorial. 
```bash
mkdir podinate-quick-start
cd podinate-quick-start
```

Now copy the following into `hello.pf` to create a super advanced hello world application.
```hcl title="hello.pf"
pod "hello-world" {
    image = "ubuntu"
    command = [ "/bin/bash", "-c", "--" ]
    arguments = [ "while true; do echo 'Hello from Podinate!'; sleep 2; done;" ]
}
```
This file creates an Ubuntu pod which will log "Hello from Podinate!" every two seconds. You can create the Ubuntu Pod by running;
```bash
podinate apply hello.pf
```
Podinate will show you exactly what it will create in the Kubenetes cluster: 
```
namespace default is up to date
pod hello-world will be  created :
apiVersion: apps/v1
kind: StatefulSet
metadata:
  creationTimestamp: null
  name: hello-world
  namespace: default
spec:
  persistentVolumeClaimRetentionPolicy:
    whenDeleted: Retain
    whenScaled: Retain
  replicas: 1
  selector:
    matchLabels:
      podinate.com/pod: hello-world
  serviceName: ""
  template:
    metadata:
      creationTimestamp: null
      labels:
        podinate.com/pod: hello-world
      name: hello-world
    spec:
      containers:
      - args:
        - while true; do echo 'Hello from Podinate!'; sleep 2; done;
        command:
        - /bin/bash
        - -c
        - --
        image: ubuntu:latest
        name: hello-world
        resources: {}
  updateStrategy: {}
status:
  availableReplicas: 0
  replicas: 0
Summary: 1 created, 0 updated, 0 deleted, 1 unchanged
Are you sure you want to apply these changes? (Y/n)
```
Podinate is showing that it will create a Kubernetes *StatefulSet*, which will in turn run the Pod we want. Podinate will create various Kubernetes *Objects* for you when you apply a PodFile, and will confirm exactly what changes will be made whenever you run `apply`. This lets you see at a glance exactly what Podinate is changing. 

Apply the `hello.pf` PodFile by pressing Y and enter. 

### List Pods
To list pods on Kubernetes, use the `kubectl` tool. It should have been installed with k3d by Brew. 
```bash
kubectl get pods 
```
Because our Ubuntu pod is the only thing running in the default namespace, the output should look like the following: 

```
NAME            READY   STATUS    RESTARTS   AGE
hello-world-0   1/1     Running   0          38s
```
### Check Pod Logs 
The Pod logs will contain the output of the program running inside the container. In this case, it's our super advanced "hello world" system. 
```bash
kubectl logs -f hello-world-0
```
You'll now be seeing "Hello from Podinate!" logged every two seconds. This is a very useful tool to see what is going on in the Pod. The `-f` means to keep following the log and show any new lines that come up. You can omit it if you just want to see the current contents of the log. 

### Run Command in Pod
Kubectl can run any command inside of our pod. This command will start a Bash shell inside of the Pod: 
```bash
kubectl exec -it hello-world-0 -- bash
```
*hacker voice* you're in. Try running some Linux commands. 

You can also try running commands directly (very useful for scripting), such as:
```bash
kubectl exec hello-world-0 -- ls /var
kubectl exec hello-world-0 -- echo "Hello world"
```
### Change the Image to Fedora
Let's say we meant to run our super advanced hello world system on a Fedora Pod. Update the `hello.pf` PodFile as follows: 
```hcl title="hello.pf"
pod "hello-world" {
    image = "fedora"
    command = [ "/bin/bash", "-c", "--" ]
    arguments = [ "while true; do echo 'Hello from Fedora!'; sleep 2; done;" ]
}
```
Now run the following command to update the Pod: 
```bash
podinate apply hello.pf
```
Podinate will show you exactly what needs to be updated:
```
namespace default is up to date
pod hello-world will be  updated 
  apiVersion: "apps/v1"
  kind: "StatefulSet"
  metadata:
    creationTimestamp: "2024-07-07T02:11:26Z"
-   generation: 1
+   generation: 2
    name: "hello-world"
    namespace: "default"
    resourceVersion: "1676"
    uid: "e6c60e53-136a-44a7-9f83-7998b88d5fa1"
  spec:
    persistentVolumeClaimRetentionPolicy:
      whenDeleted: "Retain"
      whenScaled: "Retain"
    podManagementPolicy: "OrderedReady"
    replicas: 1
    revisionHistoryLimit: 10
    selector:
      matchLabels:
        podinate.com/pod: "hello-world"
    serviceName: ""
    template:
      metadata:
        creationTimestamp:
        labels:
          podinate.com/pod: "hello-world"
        name: "hello-world"
      spec:
        containers:
          -
            args:
-             - "while true; do echo 'Hello from Podinate!'; sleep 2; done;"
+             - "while true; do echo 'Hello from Fedora!'; sleep 2; done;"
            command:
              - "/bin/bash"
              - "-c"
              - "--"
-           image: "ubuntu:latest"
+           image: "fedora:latest"
            imagePullPolicy: "Always"
            name: "hello-world"
            resources:
            terminationMessagePath: "/dev/termination-log"
            terminationMessagePolicy: "File"
        dnsPolicy: "ClusterFirst"
        restartPolicy: "Always"
        schedulerName: "default-scheduler"
        securityContext:
        terminationGracePeriodSeconds: 30
    updateStrategy:
      rollingUpdate:
        partition: 0
      type: "RollingUpdate"
  status:
    availableReplicas: 1
    collisionCount: 0
    currentReplicas: 1
    currentRevision: "hello-world-79ff6c6b6c"
    observedGeneration: 1
    readyReplicas: 1
    replicas: 1
    updateRevision: "hello-world-79ff6c6b6c"
    updatedReplicas: 1

Summary: 0 created, 1 updated, 0 deleted, 1 unchanged
Are you sure you want to apply these changes? (Y/n)
```
With a quick glance, you can check that Podinate will only change what you changed in the PodFile. Apply the configuration by pressing Y and enter.

Now watch the changes be applied: 
```bash
kubectl get pods -w
```
The `-w` means `watch`, and will show the Ubuntu Pod being replaced by a Fedora Pod. This may take a minute or two. 
```
NAME            READY   STATUS        RESTARTS   AGE
hello-world-0   1/1     Terminating   0          4m9s
hello-world-0   0/1     Terminating   0          4m32s
hello-world-0   0/1     Terminating   0          4m33s
hello-world-0   0/1     Terminating   0          4m33s
hello-world-0   0/1     Terminating   0          4m33s
hello-world-0   0/1     Pending       0          0s
hello-world-0   0/1     Pending       0          0s
hello-world-0   0/1     ContainerCreating   0          0s
hello-world-0   1/1     Running             0          11s
```
Try running `kubectl logs -f hello-world-0` again to see how the Pod's log output has changed, or running some commands in the Pod to see how the environment has changed. 

## Next Steps
This tutorial should have demonstrated how easy it is to control Kubernetes with Podinate. Here's some ideas for what to do next:

- Try changing the Pod to log the time with each "Hello from Fedora!" 
- [Set up First App](first-app.md) is part two of this tutorial, and will take you through setting up Uptime Kuma, which you can use to monitor every other app you run on your Kubernetes cluster.
