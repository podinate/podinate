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