#!/bin/bash
# Install Salt Minion
curl -L https://bootstrap.saltproject.io -o install_salt.sh
sudo sh install_salt.sh -P
echo "master: salt.podinate.com" | sudo tee /etc/salt/minion.d/master.conf
sudo systemctl restart salt-minion