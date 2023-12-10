base: 
    '*':
        - common
        - zabbix.agent
    'podinate-zabbix*':
        - zabbix.server
    'internal-*':
        - kzero.server
    'internal-1.podinate.com':
        - kzero.base
    'podinate-salt*':
        - kzero.controller

