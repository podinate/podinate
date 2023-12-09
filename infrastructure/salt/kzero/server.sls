kube_server_packages:
    pkg.installed:
        - pkgs:
            - open-iscsi
            - nfs-common

iscsid:
    service.running:
        - enable: true
        - require:
            - pkg: kube_server_packages