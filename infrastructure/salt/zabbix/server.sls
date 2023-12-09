## Warning, you will need to have set up the Zabbix repo for this to work correctly 
zabbix-pkgs:
  pkg.installed:
    - pkgs:
      - zabbix-server-pgsql
      - zabbix-frontend-php
      - php8.1-pgsql
      - zabbix-nginx-conf
      - zabbix-sql-scripts
      - zabbix-agent
      - postgresql
    
zabbix-server:
  service.running:
    - name: zabbix-server
    - enable: True
    - require:
      - pkg: zabbix-pkgs

zabbix-agent:
  service.running:
    - name: zabbix-agent
    - enable: True
    - require:
      - pkg: zabbix-pkgs

nginx:
  service.running:
    - name: nginx
    - enable: True
    - require:
      - pkg: zabbix-pkgs

php-fpm:
  service.running:
    - name: php8.1-fpm
    - enable: True
    - require:
      - pkg: zabbix-pkgs

postgresql:
  service.running:
    - name: postgresql
    - enable: True
    - require:
      - pkg: zabbix-pkgs