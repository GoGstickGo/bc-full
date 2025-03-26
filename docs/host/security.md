# User Account Security
## Create a non-root user for daily operation
```shell
sudo adduser nodeadmin
sudo usermod -aG sudo nodeadmin
```
## Enforce strong password policies
```shell
sudo apt install libpam-pwquality
sudo nano /etc/pam.d/common-password
    password requisite pam_pwquality.so retry=3 minlen=12 difok=3 ucredit=-1 lcredit=-1 dcredit=-1 ocredit=-1
```
## Disable root login
```shell
sudo passwd -l root
```
## SSH Hardening
```shell
sudo nano /etc/ssh/sshd_config
    PermitRootLogin no
    PasswordAuthentication no
    PubkeyAuthentication yes
    PermitEmptyPasswords no
    Protocol 2
    X11Forwarding no
    AllowUsers nodeadmin  # Replace with your username
    Port 2222  # Change default port
```
# System Hardening
## Firewall Configuration
```shell
sudo apt install ufw
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 2222/tcp  # SSH port you configured
sudo ufw allow 8333/tcp  # Bitcoin Core
sudo ufw allow 9735/tcp  # Lightning Network
sudo ufw enable

sudo ufw status verbose
```
## Keep the system updated
```shell
sudo apt update && sudo apt upgrade -y
sudo apt install unattended-upgrades
sudo dpkg-reconfigure unattended-upgrades
```
## Remove unnecessary services and packages
```shell
sudo apt remove --purge telnet rsh-server rsh-client
sudo apt autoremove
```
## Enable automatic security updates
```shell
sudo nano /etc/apt/apt.conf.d/50unattended-upgrades
 # Uncomment the security updates line
```
## Disable unused filesystems
```shell
echo "install cramfs /bin/true" | sudo tee -a /etc/modprobe.d/disable-filesystems.conf
echo "install freevxfs /bin/true" | sudo tee -a /etc/modprobe.d/disable-filesystems.conf
echo "install jffs2 /bin/true" | sudo tee -a /etc/modprobe.d/disable-filesystems.conf
```
# Kernel hardening
## Secure shared memory
```shell
sudo sudo nano /etc/fstab
# add tmpfs /run/shm tmpfs defaults,noexec,nosuid 0 0
```
## Adjust kernel parameters
```shell
sudo nano /etc/sysctl.conf

    kernel.randomize_va_space = 2
    net.ipv4.conf.all.rp_filter = 1
    net.ipv4.conf.default.rp_filter = 1
    net.ipv4.conf.all.accept_redirects = 0
    net.ipv6.conf.all.accept_redirects = 0
    net.ipv4.conf.all.send_redirects = 0
    net.ipv4.conf.all.accept_source_route = 0
    net.ipv6.conf.all.accept_source_route = 0
```
```shell
sudo sysctl -p
# Apply changes
```
# Additional Security Measures
## fail2ban
```shell
sudo apt install fail2ban
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
sudo nano /etc/fail2ban/jail.local
    bantime = 3600
    findtime = 600
    maxretry = 3
```
## Enable process accounting
```shell
sudo apt install acct
sudo touch /var/log/wtmp
sudo systemctl enable acct
sudo systemctl start acct
```
### Read wtmp logs
```shell
# Show only the last 20 entries
last -20

# Show logins by a specific user
last nodeadmin

# Show system reboots
last reboot

# Show login duration (how long users were logged in)
last -F

# Specify a different wtmp file if needed
last -f /var/log/wtmp.1
```
## Security monitoring tools
```shell
sudo apt install rkhunter lynis
```
## Configure system auditing
```shell
sudo apt install auditd
sudo systemctl enable auditd
sudo systemctl start auditd
```
