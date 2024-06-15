# Cluster Commands
This directory contains all the commands related to Podinate's built-in cluster management. These commands are pulled in from the autok3s package from ranchercn. https://github.com/cnrancher/autok3s 

## Autok3s 
Autok3s is a package maintained by the CN division of SUSE Rancher. I haven't been able to find any overall documentation, only the [documentation of the providers here](https://github.com/cnrancher/autok3s/tree/master/docs/i18n/en_us)

### GUI
There is a fairly straightforward GUI to manage Kubernetes clusters which is accessed by just running the `autok3s` command. 

### Commands
Autok3s uses Golang Cobra package just like Podinate so it's convenient we can simply pull in the command wholesale and patch into the PostRun on the commands. 

I haven't found very good documentation of the usage of the commands so I'm going to put them here for now. 

The CLI itself isn't very well documented either. For example, there is nothing attached to the Cobra built-in help for the various commands listed below.

#### Create
Creates a K3s cluster using the provider of your choice. 
Available providers are documented in the Readme of the autok3s project, but these are the most relevant:
- K3d (local cluster, for Podinate users for testing and development)
- AWS 
- Google Cloud
- Bare Metal (SSHes into the machines and sets up the K3s installer)

```bash
p cluster create --provider k3d --name local --master 1
``` 
This will create a K3d cluster called local with one node. When Podinate is installed I'd like to create this cluster by default. 

#### List 
Pretty straightforward, just shows a list of all the clusters that k3s has created. 
```bash
p cluster list
  NAME   REGION  PROVIDER  STATUS   MASTERS  WORKERS    VERSION     ISHAMODE  DATASTORETYPE  
  local          k3d       Running  1        0        v1.23.8+k3s1  false                  
```

#### Describe
Prints out some YAML describing a given cluster in a bit more detail.
```yaml
$ p cluster describe -n local
Name: local
Provider: k3d
Region: 
Zone: 
Master: 1
Worker: 0
IsHAMode: false
Status: Running
Version: v1.23.8+k3s1
Nodes:
  - internal-ip: []
    external-ip: []
    instance-status: running
    instance-id: k3d-local-server-0
    roles: control-plane,master
    status: Ready
    hostname: k3d-local-server-0
    container-runtime: containerd://1.5.13-k3s1
    version: v1.23.8+k3s1
```
I'd kind of prefer `p cluster describe local` but this will do for now. 

#### Delete
Deletes the cluster. 
```bash
$ p cluster delete -p k3d -n local
? [k3d] are you sure to delete cluster local Yes
time="2024-06-03T15:52:25+12:00" level=info msg="[k3d] begin to delete cluster local..."
time="2024-06-03T15:52:25+12:00" level=info msg="Deleting cluster 'local'"
time="2024-06-03T15:52:26+12:00" level=info msg="Deleting cluster network 'k3d-local'"
time="2024-06-03T15:52:26+12:00" level=info msg="Deleting 1 attached volumes..."
time="2024-06-03T15:52:26+12:00" level=info msg="[k3d] successfully delete cluster local"
time="2024-06-03T15:52:27+12:00" level=info msg="[k3d] successfully deleted cluster local"
```