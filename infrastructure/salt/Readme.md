# Control Infrastructure
This folder has the Salt config to set up the salt master, zabbix, and the API and user Kubernetes clusters.

## Purpose
Salt's role is to bootstrap the underlying VMs, install Kubernetes (k3s) and dependencies like Github and Longhorn, then allow other scripts to deploy the API backend and any software that needs to be installed on the runner clusters.

## Bootstrapping
Create two new VMs with Ubuntu. One will be the salt master, the other will be the Zabbix server.

- On the salt master, run the standard salt bootstrap script.
### Zabbix
- On the Zabbix server, run the `bootstrap_minion.sh` file. 
- On the Zabbix server, run the handful of commands at the top of the Zabbix install page to set up the repo.
```bash
wget https://repo.zabbix.com/zabbix/6.4/ubuntu/pool/main/z/zabbix-release/zabbix-release_6.4-1+ubuntu22.04_all.deb
dpkg -i zabbix-release_6.4-1+ubuntu22.04_all.deb
apt update 
```

