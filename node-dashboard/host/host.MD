```bash
[Unit]
Description=Full Node Monitor
After=network.target

[Service]
Type=Simple
User=${USER}
WorkingDirectory=/home/${USER}/${FOLDERPATH}
ExecStart=/home/${USER}/${FOLDERPATH}/node-dashboard -rpcuser=${RPCUSER} -rpcpass=${PW}-rpchost=localhost:${RPCIP} -listen=localhost:${HTTPSERVERIP}
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```