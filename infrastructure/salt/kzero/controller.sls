# k0s saltstack formula
# Installs k0sctl and kubectl, then uses k0sctl to install a k0s cluster

# Install k0sctl and kubectl
/usr/local/bin/k0sctl:
  file:
    - managed
    - source: https://github.com/k0sproject/k0sctl/releases/download/v0.16.0/k0sctl-linux-x64
    - source_hash: b21eb9edc90180d3bfb3413b82efac3811989f85d5c1f8e54750bf7b2144fbc6
    - user: root
    - group: root
    - mode: 0755

/usr/local/bin/kubectl:
  file:
    - managed
    - source: https://dl.k8s.io/release/v1.28.4/bin/linux/amd64/kubectl
    - source_hash: 893c92053adea6edbbd4e959c871f5c21edce416988f968bec565d115383f7b8
    - user: root
    - group: root
    - mode: 0755

/usr/local/bin/k:
  file.symlink:
    - target: /usr/local/bin/kubectl
    - require:
        - file: /usr/local/bin/kubectl

# Use k0sctl to set up the cluster
install_cluster:
  cmd.run:
    - name: k0sctl apply --config /srv/salt/kzero/k0sctl.yaml
    - user: ubuntu
    - group: ubuntu
    - cwd: /home/ubuntu
    - shell: /bin/bash
    - require:
        - file: /usr/local/bin/k0sctl

# Grab the kubeconfig and install it 
/home/ubuntu/.kube:
  file.directory:
    - user: ubuntu
    - group: ubuntu
    - mode: 0755

export_kubeconfig:
  cmd.run:
    - name: k0sctl kubeconfig --config /srv/salt/kzero/k0sctl.yaml > /home/ubuntu/.kube/config
    - user: ubuntu
    - group: ubuntu
    - cwd: /home/ubuntu
    - shell: /bin/bash
    - require:
        - cmd: install_cluster
        - file: /home/ubuntu/.kube