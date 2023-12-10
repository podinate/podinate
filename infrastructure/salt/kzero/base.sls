## Base install - set up Longhorn and other common tools for all nodes

install_helm:
    pkg.installed:
      - pkgs:
        - curl
        - gnupg2
    cmd.run:
      - cwd: /root
      - names:
        - curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
        - chmod 700 ./get_helm.sh
        - ./get_helm.sh


longhorn_helm_repo:
    helm.repo_managed:
      - present:
        - name: longhorn
          url: https://charts.longhorn.io
      - require:
          - pkg: install_helm

update_helm_repos:
    helm.repo_updated:
      - name: longhorn
      - require:
        - pkg: install_helm
        - helm: longhorn_helm_repo

longhorn_install:
    helm.release_present:
      - name: longhorn
      - namespace: longhorn-system
      - chart: longhorn/longhorn
      - flags:
        - dry-run
      - values: /srv/salt/kzero/files/longhorn.values.yaml
      - require:
        - pkg: install_helm
        - helm: update_helm_repos