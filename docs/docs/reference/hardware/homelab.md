# Podinate for Homelab
Podinate is perfect for running apps and containers in a homelab environment. It allows you to quickly install many widely used applications and quickly deploy and prototype your own projects. 

Podinate runs an example homelab setup, and this page documents the setup of our example system. We will mention alternatives where possible. 

## Hardware
Your hardware requirements depend on what you want to run. Podinate can run on machines with a single core and 2GB of ram, but will only have enough headroom to run one or two applications in such a constrained setup.

{{ read_csv('./homelab_hardware.csv') }}

For a homelab environment it's recommended to stick with a single node, as this will make storage and other management a lot easier. 

## Storage
For home lab storage we keep things easy. That means using the [LVM](../storage/LVM.md) or [ZFS](../storage/ZFS.md).

### Single Disk 
If you have a single large disk you wish to use, we recommend setting it up as an LVM volume group and using the [LVM](../storage/LVM.md) configuration to connect it to Podinate. In our example homelab setup we have a single 1TB NVMe SSD that serves as high-performance storage for our pods. 

### Multi-disk / RAID
If you want to connect multiple disks together to Podinate, we recommend using the battle-tested [ZFS](../storage/ZFS.md).

## Workloads
Some of the most common workloads for a homelab include:

- Personal Cloud
    - NextCloud
- Media server 
    - Plex
    - Emby / Jellyfin 
- Productivity Software
    - Penpot
    - OnlyOffice Server

## Install
The installation of the server for the Podinate reference Homelab went as follows:

- Install Ubuntu Server.
    - [Download Ubuntu Server](https://ubuntu.com/download/server) ISO
    - Burn the ISO to a USB flash drive using [Balena Etcher](https://etcher.balena.io/)
    - Start the machine, pressing the `esc` key on the keyboard to interrupt boot. 
    - In the menu that appears, select the USB flash drive.
- Perform Ubuntu Install
    - Proceed through the first few steps of the installer 
    - When you get to the partitioning step, select use entire disk
    - Enable [full disk encryption](../../storage/lvm#encryption) if you can, 

