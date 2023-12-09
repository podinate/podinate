install_zabbix_agent:
   pkg.installed:
      - name: zabbix-agent

/etc/zabbix/zabbix_agentd.conf:
   file:
      - managed
      - source: salt://zabbix/zabbix_agentd.conf
      - template: jinja
      - require:
        - pkg: zabbix-agent

run_agent:
   service.running:
      - name: zabbix-agent.service
      - enable: true
      - restart: true
      - watch:
        - file: /etc/zabbix/zabbix_agentd.conf
      - require:
        - pkg: zabbix-agent
        - file: /etc/zabbix/zabbix_agentd.conf
        - file: /var/log/zabbix_agentd.log

/var/log/zabbix_agentd.log:
   file:
      - managed
      - user: zabbix
      - group: zabbix
      - mode: 644
      - require:
        - pkg: zabbix-agent
