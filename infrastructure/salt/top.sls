base: 
    '*':
        - common
        - zabbix.agent
    'podinate-zabbix*':
        - zabbix.server
    'internal-*':
        - kzero.server
    'podinate-salt*':
        - kzero.controller

