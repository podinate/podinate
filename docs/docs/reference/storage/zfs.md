# ZFS with Podinate
ZFS (Zettabyte File System) is a massively scalable filesystem originally developed by Sun for their Solaris operating system. Development of ZFS now mostly centers around the open source OpenZFS project which continued work on the filesystem after Sun's acquisition. 

Podinate can use ZFS based storage through the [OpenEBS ZFS Provisioner](https://github.com/openebs/zfs-localpv)

## Terms
ZFS has a couple of terms you should be familiar with: 

- **pool** - A ZFS pool is a set of one or more disks that can be used to create *datasets*. By default a pool is also a dataset. 
- **Dataset** - A Dataset is similar to a partition, and can be used to store files / mounted. They are usually referred to like `poolname/setname` by the OpenZFS command. They can also be nested, such as `pool/podinate/pvc-abc123`

## Compression
ZFS supports native compression. It is generally recommended to turn on in nearly all circumstances, as the compression algorithm is likely faster than all but the highest end NVMe SSD drives. 

## Deduplication
Deduplication uses a ton of RAM to detect blocks of identical data on the disks. Generally this is not recommended, unless you have a use case that requires storing a ton of duplicated data for some reason. Even then, the deduplication only works if duplicate blocks line up exactly on the disk. 

## RAID
ZFS has built in RAID support. You may be familiar with RAID levels, such as RAID0 or RAID 6. ZFS changes the names of the RAID levels slightly. 

{{ read_csv('./ZFS_raid_levels.csv') }}

## Creating a ZFS Pool
Setting up a ZFS pool for Podinate is a pretty quick process. If you want to use just a single disk, Podinte recommends using [LVM](../lvm) instead.

1. Install ZFS on Ubuntu
    ```bash
    sudo apt-get install zfsutils
    ```
1. Find your disks. The following commands will help you find the /dev/sd* IDs of your disks. 
    ```bash
    sudo fdisk -l 
    ls -l /dev/disk/by-id
    ```
1. Option 1: Create an unencrypted *pool*. 
    Create pool with two disks of redundancy and compression enabled. See [RAID](#raid) for the levels. To create a pool of one disk or a striped set, omit the raid level. 
    ```bash
    zpool create -O compression=enable vault raidz2 /dev/sda /dev/sdb /dev/sdc /dev/sdd /dev/sde 
    ```
    Create with single disk / no redundancy (NOT recommended):
    ```bash
    zpool create vault /dev/sda
    ```
1. Option 2: Create a pool encrypted by a keyfile. 
    If you use a keyfile on disk, **the filesystem storing the key must be encrypted through some other means**, for example [LVM](../lvm#encryption) full disk encryption. 
    ```bash
    # Create the encryption key
    head -c 32 /dev/random > /podinate/key/zfs

    # Create the pool 
    zpool create \
    -O compression=lz4 \
    -O encryption=aes-256-gcm \
    -O keyformat=raw \
    -O keylocation=file:///podinate/key/zfs \
    vault raidz2 \
    /dev/sda \ # Make sure to change to your actual disks
    /dev/sdb \
    /dev/sdc \
    /dev/sdd \
    /dev/sde 
    ```
    ZFS will remember where the key is stored and automatically mount the pool at boot. 
1. Create a *Dataset* for Podinate:
    Datasets reside within a pool and are a way to divide up your ZFS pool. It's a good idea to create a separate dataset for Podinate, so you can use the pool for other things like [Incus](../../software/incus) Virtual Machines later.
    ```bash
    zfs create vault/podinate 
    ```
1. (Encrypted) Unlock Encrypted Dataset on Boot:
    ZFS won't load the encryption key into our Podinate Dataset by default for some reason, so we have to make it. 
    ```bash
    nano /etc/systemd/system/zfs-load-keys.service
    ```
    Copy in the following contents: 
    ```systemd
    [Unit]
    Description=Load all ZFS encryption keys
    DefaultDependencies=no
    Before=zfs-mount.service
    After=zfs-import.target
    Requires=zfs-import.target

    [Service]
    Type=oneshot
    RemainAfterExit=yes
    ExecStart=/usr/sbin/zfs load-key -a

    [Install]
    WantedBy=zfs-mount.service
    ```
    And enable the service by running: 
    ```bash
    systemctl enable zfs-load-keys
    ```
    All dataset keys will be loaded automatically on boot. 

## Conecting Pool to Podinate 
Once we have a Dataset to use for Podinate, the next step is to set up the ZFS volume provisioner in our Kubernetes cluster.

!!! warn 
    You cannot edit a Kubernetes StorageClass after creating it. Check the [ZFS OpenEBS Docs](https://github.com/openebs/zfs-localpv/blob/develop/docs/storageclasses.md) for parameters. If you need to change something after volumes already exist, consider creating a new StorageClass and using it alongside the existing class. 


1. Install the OpenEBS ZFS Provisioner
    ```bash
    kubectl apply -f https://openebs.github.io/charts/zfs-operator.yaml
    ```

1. Create a storage class to use the dataset

    ```yaml
    apiVersion: storage.k8s.io/v1
    kind: StorageClass
    metadata:
        name: bulk
    parameters:
        compression: "on"
        dedup: "off"
        fstype: "zfs"
        poolname: "vault/podinate"
        shared: "yes"
        thinprovision: "yes" 
    provisioner: zfs.csi.openebs.io
    volumeBindingMode: WaitForFirstConsumer
    ```
