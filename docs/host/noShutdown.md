# Set Up Automatic Restart After Power Failure
## Configure your system to automatically power on after a power outage:
```shell
sudo nano /etc/default/grub
# Find the line with GRUB_CMDLINE_LINUX_DEFAULT and add consoleblank=0
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash consoleblank=0"
# Update grub:
sudo update-grub
```
## Configure Kernel to Avoid System Freezes
```shell
sudo nano /etc/sysctl.conf

Improve system stability
kernel.panic = 10
kernel.panic_on_oops = 1
vm.swappiness = 10

sudo sysctl -p
```
## Disable Automatic Updates That Require Reboots
```shell
sudo nano /etc/apt/apt.conf.d/50unattended-upgrades

Unattended-Upgrade::Automatic-Reboot "false";
```
## Use a Process Manager
```shell
/etc/systemd/system/bitcoind.service

[Unit]
Description=Bitcoin daemon
After=network.target

[Service]
User=nodeadmin
Group=nodeadmin
Type=forking
ExecStart=/usr/local/bin/bitcoind
ExecStop=/usr/local/bin/bitcoin-cli stop
Restart=always
RestartSec=30
TimeoutStopSec=30min

[Install]
WantedBy=multi-user.target
```

### Reload systemd configuration
```shell
sudo systemctl daemon-reload
```
