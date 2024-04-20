# LVM on Kubernetes / Podinate
LVM stands for Logical Volume Manager. It is a framework that provides simplified disk and storage management for Linux. It abstracts away things like hard disks, partitions and filesystems. 

Podinate can use LVM based storage through the [OpenEBS LVM Provisioner](https://github.com/openebs/lvm-localpv)

## Concepts
LVM is an abstraction layer between partitions and the physical disks on which they reside. This makes it easier to do a lot of operations as they're now done through LVM instead of physically on disk.

LVM lets you group together multiple disks and address them as one, but it's important to remember **it does not provide any redundancy in this mode**. For redundancy / RAID, you need to use something like [ZFS](../ZFS).

LVM uses a few terms that you should be familiar with: 

- **Physical Volume** - A physical device, usable by LVM for storage. These can be a hard disk, a partition on a hard disk, or something more exotic like a SAN array. 
- **Volume Group** - A group of Physical Volumes that functions as a single logical device to provision Logical Volumes on. 
- **Logical Volumes** - The LVM equivalent of partitions. This is where filesystems are created and files are stored. 

To create partitions on LVM, you create a *Volume Group* that spans one or more *Phyisical Volumes*, then create *Logical Volumes* on the *Volume Group*.

## Encryption
LVM supports encrypting the entire Volume Group. If you're installing the recommended Ubuntu Server setup, it will ask if you want to encrypt the entire disk during the install. This is highly recommended but does require some additional setup. 
!!! tip 
    An encrypted "Root Volume Group" will require you to enter the decryption password physically into the machine before it will boot, unless you set up [Dropbear](../software/dropbear.md) to enable remote unlocking. 

## Snapshots
LVM enables taking a snapshot of a *Logical Volume*. This means that regardless of the filesystem used, you can get a consistent snapshot and form of backup.

## Create a Volume Group for Use With Podinate
The [OpenEBS LVM Provisioner](https://github.com/openebs/lvm-localpv) can use an LVM *Volume Group* to provide storage to Pods in a Podinate cluster. It will create *Logical Volumes* as Pods need them. 

1. **Find the Drive to Use**

    The following commands will help to find the drive:
    ```bash
    lsblk
    ```
    This will show a tree of all the drives installed in the system. The drive should be identifiable from its size. 
    Try this command if you're not sure: 
    ```bash
    ls -l /dev/disk/by-id
    ```
    This will list the devices by the ID the manufacturer gave them, with the location of the disk. 

1. **Erase any existing partitions**

    ```bash
    sudo parted /dev/sda rm 1
    ```
    This will delete the first partition (which Linux refers to as `/dev/sda1`) of the drive. 

1. (Optional) **Encrypt the Drive with Luks**

    This will set up full disk encryption on the disk, so nobody can read it if the disk is removed from the system. 
    ```bash
    sudo cryptsetup luksFormat /dev/sda
    ```
    Enter a passphrase for the disk. If [Dropbear](../software/dropbear.md) is set up, the disk can be unlocked remotely by an administrator on boot using this passphrase.

1. (Optional) **Automatically Unlock Disk On Boot**

    !!!tip
        Make sure the key file generated is stored on a disk that is itself encrypted through some other means. If full disk encryption was used on the host Ubuntu OS, the key file could be stored there. 
    Create a randomly generated key for use with the luks Drive
    ```bash
    dd bs=512 count=4 if=/dev/random of=/etc/crypt.key
    sudo cryptsetup luksAddKey /dev/sda /etc/crypt.key
    sudo cryptsetup luksOpen /dev/sda/ ssd-disk
    ```
    The final `ssd-disk` is an ID you want to refer to the open device by, for example `fast` or `ssd`.

    To unlock the disk on boot the disk has to be added to the `crypttab`. First get the UUID of the disk:
    ```bash
    lsblk -f 
    ```
    Copy the random string in the fifth column next to the encrypted file. 
    ```bash
    nano /etc/crypttab 
    ```
    Add the following line to the bottom of the file:
    ```bash
    ssd-disk UUID=your-random-id /etc/crypt.key luks
    ```

1. **Create the LVM Volume Group**
    ```bash
    # If using an encrypted volume
    sudo vgcreate podinate /dev/mapper/ssd-disk
    # If using an unencrypted disk
    sudo vgcreate podinate /dev/sda
    ```

1. **Install OpenEBS LVM Provisioner**
    
    ```bash
    kubectl apply -f https://openebs.github.io/charts/lvm-operator.yaml
    ```
    Check all the pods are running: 
    ```bash
    kubectl get pods -n openebs -l role=openebs-lvm
    ```
    The output should look like:
    ```
    NAME                                              READY   STATUS    RESTARTS   AGE
    openebs-lvm-localpv-controller-7b6d6b4665-fk78q   5/5     Running   0          11m
    openebs-lvm-localpv-node-mcch4                    2/2     Running   0          11m
    ```
    Create a storageclass to use the provisioner: 
    ```yaml
    cat sc.yaml 

    apiVersion: storage.k8s.io/v1
    kind: StorageClass
    metadata:
      name: ssd-lvm
    parameters:
      storage: "lvm"
      volgroup: "podinate"
    provisioner: local.csi.openebs.io
    ```

1. **Use the Storage Class**
    Create a Shared Volume with the storage class: 
    ```hcl
    shared_volume "lvm" {
        size = 2
        class = "ssd-lvm"
    }
    ```
    Add a non-shared volume to a Pod with the storage class:
    ```hcl
    pod "ubuntu" {
        image = "ubuntu"
        command = [ "/bin/bash", "-c", "--" ]
        arguments = [ "while true; do sleep 120; done;" ]
        volume "lvm" {
            size = 2
            path = "/data"
            class = "ssd-lvm"
        }
    }
    ```

## References
- [LVM Official Website](https://sourceware.org/lvm2/)
- [Arch Wiki - LVM](https://wiki.archlinux.org/title/LVM) 
- [Digital Ocean - Introduction to LVM](https://www.digitalocean.com/community/tutorials/an-introduction-to-lvm-concepts-terminology-and-operations)
- [OpenEBS LVM Provisioner](https://github.com/openebs/lvm-localpv)
