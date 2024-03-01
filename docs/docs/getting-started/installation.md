# Server Install

Podinate is a Kubernetes controller that provides a simple interface to an underlying Kubernetes cluster. This lets you spin up apps and containers quickly without babysitting the underlying Kubernetes or writing acres of YAML. 

This documentation will guide you through setting up a dedicated Podinate cluster. If you want to install a development cluster on your laptop, that will be a different process. 

## System Requirements
We recommend using a dedicated VM or hardware node for your Podinate server. Podinate can run in as little as 2GB of ram for edge compute situations, but the number of pods you can run will be limited by how much ram you can provide. 

Podinate supports multiple storage types. For example you might want to have a separate NVMe SSD for pods such as database servers, and have an array of large HDDs for storing bulk files, depending on your needs.

- Clean Ubuntu 22.04 install 
- CPU: 2 Cores (4+ Recommended)
- RAM: 8GB (16+ Recommended)
- Disk: 128GB SSD Minimum
    - Recommended:
        - 128GB Ubuntu install disk 
        - 256GB+ SSD for Pod storage
        - 1TB+ HDD for bulk storage if desired

## 1. Install Ubuntu
[Download Ubuntu Server](https://ubuntu.com/download/server) from the official website, burn it to a USB using [Etcher](https://etcher.balena.io/), and boot the machine from the USB drive. If you're setting up a virtual machine, mount the downloaded ISO file into the virtual machine using your virtual machine manager. If you have an existing Ubuntu server and you'd just like to try Podinate, we recommend [Incus](https://linuxcontainers.org/incus/) to create and manage virtual machines. 

Follow the installer steps until you get to partitioning. If you select to use automatic partitioning, the installer will probably set up a 100GB `/` root partition. We recommend creating a single partition using the rest of your disk, and setting it to be mounted at `/var` as Podinate will store all your pod data under there. 

## 2. Install K3S
K3S is a minimal distribution of Kubernetes that we recommend for Podinate usage. If your company has a requirement for a different version of Kubernetes, advanced users could try skipping this step and installing on your existing environment.

[K3s Installation Docs](https://docs.k3s.io/quick-start)

```bash
curl -sfL https://get.k3s.io | sh -
```

## 3. Set Up Disks 
If you installed Ubuntu on a single disk with a separate `/` and `/var` partition, and don't want to add additional storage, you can skip this step. 

K3s' default storage location is `/var/lib/rancher/k3s/storage`. We recommend mounting whatever disk or array of disks you want to be your default storage location for pods in this location.

Use the following command to see your disks:
```bash
lsblk
```
This will show you all the disks attached to the system, and you should be able to see which one you want to use as Pod storage. Make a note of its /dev/sdax number. Then use the following command to find the disk's ID.
```bash
ls -l /dev/disk/by-id/
``` 
This will show you the IDs of the disks and partitions, and which sdx device they point to. Use this to figure out which disk you want to mount. The following command will mount the disk at the correct location and ensure it will be mounted there after reboot. 
```bash
mkdir /var/lib/rancher/k3s/storage
mount /dev/disk/by-id/<some-id> /var/lib/rancher/k3s/storage
echo "/dev/disk/by-id/<some-id> /var/lib/rancher/k3s/storage <fstype> defaults 0 0" >> /etc/fstab
```
