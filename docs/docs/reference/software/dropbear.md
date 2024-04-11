# Dropbear
[Dropbear](https://matt.ucc.asn.au/dropbear/dropbear.html) is a small, cross-platform SSH server and client implementation designed for embedded systems. It can do many things, but the use-case most useful to Podinate is to use it to perform remote unlocking of a root partition. 

## Unlocking Encrypted Disk
These instructions apply to Ubuntu, the recommended server OS for Podinate. Dropbear unlocks encrypted disks by being installed to the [initramfs](https://wiki.gentoo.org/wiki/Initramfs), which is a small system that is loaded before Linux, its job is to mount required filesystems and boot Linux. 

### Installing Dropbear
1. Log in to the server
1. Switch to root user
    ```bash
    sudo -i 
    ```
1. Update the system
    ```bash
    apt-get update 
    apt-get upgrade
    ```
1.  Install the Dropbear package
    ```bash
    apt-get install dropbear-initramfs 
    ```
1. Configure Dropbear server options:
    ```bash
    nano /etc/dropbear/initramfs/dropbear.conf
    ```
    Paste the following line at the end of the file:
    ```
    DROPBEAR_OPTIONS="-I 180 -j -k -p 2222 -s -c cryptroot-unlock"
    ```
    The meaning of these options is: 
    - `-I 180` - Set an inactivity timeout of 180 seconds.
    - `-j` - Disable SSH local port forwarding
    - `-k` - Disable remote SSH port forwarding too
    - `-p 2222` - Listen on pot 2222. If the default port 22 is used, a warning will show on connection because the Dropbear and host SSH server keys are different. 
    - `-s` - Disable password-less login
    - `-c cryptroot-unlock` - Dropbear will ignore any command sent by the user and run cryptroot-unlock instead. This means nothing else can be done in a Dropbear session.
1. Configure static IP (optional)
    ```bash
    nano /etc/initramfs-tools/initramfs.conf
    ```
    Add the following line to the end of the file:
    ```
    IP=192.168.1.100::192.168.1.1:255.255.255.0:dropbear
    ```
    The options are as follows:
    - `192.168.1.100` - The static IP to set
    - `192.168.1.1` - The network gateway to use
    - `255.255.255.0` - The subnet mask to use
    - `dropbear` - The hostname to set
1. Add your SSH public key to dropbear:
    ```bash
    cat /home/ubuntu/.ssh/authorized_keys >> /etc/dropbear/initramfs/authorized_keys
    # OR 
    nano /etc/dropbear/initramfs/authorized_keys
    # Paste your key
    ```
1. Update the initramfs to enable our Dropbear configuration
    ```bash
    update-initramfs -u
    ```
#### Test the Installation
1. To test the instal, reboot the server and try to unlock it with Dropbear:
    ```
    sudo reboot
    ```
1. On the machine you want to unlock the disk from, run:
    ```bash
    ssh root@<ip> -p 2222
    ```
    Accept the fingerprint shown, and you will immediately be prompted for the unlock password. The output will look like:
    ```
    Please unlock disk dm_crypt-0: 
    cryptsetup: dm_crypt-0 set up successfully

    ```
    And the session will be closed.

## See Also 
- [Official Website](https://matt.ucc.asn.au/dropbear/dropbear.html)
- [CyberCiti Tutorial](https://www.cyberciti.biz/security/how-to-unlock-luks-using-dropbear-ssh-keys-remotely-in-linux/)