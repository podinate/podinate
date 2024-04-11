# LVM
LVM stands for Logical Volume Manager. It is a framework that provides logical volume management for Linux. It abstracts away things like hard disk, partitions and filesystems. 

Podinate can use LVM based storage through the [OpenEBS LVM Provisioner](https://github.com/openebs/lvm-localpv)

## Concepts
LVM is an abstraction layer between partitions and the physical disks on which they reside. This makes it easier to do a lot of operations as they're now done through LVM instead of physically on disk. It also enables nice to have features like full disk encryption by using LVM as the encryption layer.

LVM lets you group together multiple disks and address them as one, but it's important to remember **it does not provide any redundancy in this mode**. For redundancy / RAID, you need to use something like [ZFS](ZFS).

LVM uses a few terms that you should be familiar with: 

- **Physical Volume** - A physical device, usable by LVM for storage. These can be a hard disk, a partition on a hard disk, or something more exotic like a SAN array. 
- **Volume Group** - A group of Physical Volumes that functions as a single logical device to provision Logical Volumes on. 
- **Logical Volumes** - The LVM equivalent of partitions. This is where filesystems are created and files are stored. 

To create partitions on LVM, you create a *Volume Group* that spans one or more *Phyisical Volumes*, then create *Logical Volumes* on the *Volume Group*.

## Encryption
LVM supports encrypting the entire Volume Group. If you're installing the recommended Ubuntu Server setup, it will ask if you want to encrypt the Volume Group automatically. This is highly recommended but does require some additional setup. 

!!! tip 
    An encrypted "root Volume Group" will require you to enter the decryption password physically into the machine before it will boot, unless you set up [Dropbear](../../software/dropbear) to enable remote unlocking. 

## Snapshots
LVM enables taking a snapshot of a *Logical Volume*. This means that regardless of the filesystem used, you can get a consistent snapshot and form of backup. 

## References
- [LVM Official Website](https://sourceware.org/lvm2/)
- [Arch Wiki - LVM](https://wiki.archlinux.org/title/LVM) 
- [Digital Ocean - Introduction to LVM](https://www.digitalocean.com/community/tutorials/an-introduction-to-lvm-concepts-terminology-and-operations)
