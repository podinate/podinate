# Install 
This document serves as documentation of the installation procedure. It serves as notes for when we are writing the installer script later but allows for a manual install if you're brave.

This document follows the installation steps for a Prod setup. 


## System requirements: 
- Clean Ubuntu 22.04 install
- CPU: 2 Cores (4 Recommended)
- RAM: 8GB (16+ Recommended)
- Disk:
  - Recommended:
    - 32GB+ disk for ubuntu install
    - 128GB+ SSD for Pod storage
    - 1TB+ HDD for bulk storage if desired

- Install Ubuntu base OS
- Install useful packages
    ```bash
    apt install curl nano lvm2 tmux
    ```
- Install k3s 

```bash
curl -sfL https://get.k3s.io | sh -
``` 
https://docs.k3s.io/quick-start

- Mount what you want to be your default storage at `/var/lib/rancher/k3s/storage`

```bash
mkdir /var/lib/rancher/k3s/storage
mount /dev/disk/by-id/some-id /var/lib/rancher/k3s/storage
echo "/dev/disk/by-id/some-id /var/lib/rancher/k3s/storage <fstype> defaults 0 0" >> /etc/fstab
```

- Check all pods are running or completed in K3s

    ```bash
    kubectl get pods -A
    ```
- Create Podinate namespace 

    ```bash
    kubectl create namespace podinate
    ```
- Apply the postgres yaml from `kubernetes/masterdb-postgres.yaml`, and check pod is ready. For the installer we need to find a way to generate a random password here. 

```bash
kubectl apply -f masterdb-postgres.yaml
kubectl -n podinate get pods
```

- Create the podinate database tables using the sql file in `database/masterdb.sql`. Note in the future I'd like to switch to using [Goerd](https://github.com/covrom/goerd) instead of raw sql for schema management, since at the moment schema migrations are impossible. 
```bash
bash -c "kubectl -n podinate exec -it postgres-0 -- psql 'postgresql://postgres:\$\$(kubectl -n podinate get secret masterdb-secret -o jsonpath='{.data.superUserPassword}' | base64 --decode ; echo)@localhost/podinate'" < masterdb.sql
```


# Install Steps
Okay ignore what's above here's just the raw steps. 

Plan is for this to be a bash script designed to run on Ubuntu. When it comes time to install it on another platform (eg, k3d for local dev clusters), we can simply create the cluster using exec, and then create an Ubuntu Job in K8s with the relevant permissions to run this script. 

- Checks
    - Is the cluster up? 
- Ask for user email
- Install certbot
    - Kubectl apply -f raw.github.com/.../some.yaml
    - Create a cluster issuer for Let's Encrypt
- Create the Podinate Postgres instance
    - Kubectl apply -f 
- Run a Job to migrate the cluster to the latest schema
    - Using Goerd
- Install the Podinate controller
    - Kubectl apply -f from github
- Check connection to the Podinate controller

- Most of the above is just yaml wrangling, this is new:
- Create initial admin user
- Create credential for admin user
- Use credential to create default account
- Present the login information to the user (base64-encoded?)

# Demo Steps
## Introduction
- `podinate get projects`
    - Show one existing project
- `podinate -p tunnel get pods`
    - Show one pod, cloudflare tunnel
- `podinate -p tunnel logs cloudflare-tunnel`
    - Show that Cloudflare tunnel is successfully running

## Hello world 
- `cat hello-world.pod`
    - Show a project called hello-world
    - Show a pod called Ubuntu
- `podinate apply hello-world.pod`
    - Project is successfully created
    - Pod is successfully created
- `podinate -p hello-world get pods`
    - Shows the running ubuntu pod 
- `podinate -p hello-world shell hello-world`
    - Get a shell on the hello-world pod
    - Demonstrate a working Ubuntu pod
- `podinate delete hello-world.pod`
    - Pod is successfully deleted
    - Project is successfully deleted
- `podinate get projects`
    - Shows the project is deleted

## Network access
- `cat wordpress.pod`
    - Show a new project 
    - Show the two pods
    - Show the environment variables
    - Show how the pods have services
- `podinate apply wordpress.pod`
    - Project is created
    - Pods are created
    - Services are created
    - Volumes are created
- `podinate -p blog get pods`
    - Show that mariadb and wordpress are up
- `open blog.podinate.app`
    - Show working wordpress website

What's needed:

- ~~Add services to podi~~
- ~~Add volumes to podi~~


- ~~Get logs -f to work~~
- Get exec -it to work (uses above)
- Add optional command = to pod 