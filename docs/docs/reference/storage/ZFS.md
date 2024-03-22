# ZFS with Podinate
ZFS (Zettabyte File System) is a massively scalable filesystem originally developed by Sun for their Solaris operating system. Development of ZFS now mostly centers around the open source OpenZFS project which continued work on the filesystem after Sun's acquisition. 

Podinate can use ZFS based storage through the [OpenEBS ZFS Provisioner](https://github.com/openebs/zfs-localpv)


## Compression
ZFS supports native compression. It is generally recommended to turn on in nearly all circumstances, as the compression algorithm is likely faster than all but the highest end NVMe SSD drives. 

## Deduplication
Deduplication uses a ton of RAM to detect blocks of identical data on the disks. Generally this is not recommended, unless you have a use case that requires storing a ton of duplicated data for some reason. 

## RAID
ZFS has built in RAID support. You may be familiar with RAID levels, such as RAID0 or RAID 6. ZFS changes the names of the RAID levels slightly. 

{{ read_csv('./ZFS_raid_levels.csv') }}