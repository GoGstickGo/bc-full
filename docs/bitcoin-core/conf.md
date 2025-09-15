```shell
# Basic Settings
server=1
daemon=1
txindex=1

# Set the data directory to your external SSD
datadir=/mnt/bitcoin

# Network Settings
listen=1

# Performance (16GB RAM system)
dbcache=4000
maxmempool=500

# Initial Block Download Optimization
blocksonly=1

# STRICT SPAM/OP_RETURN FILTERING
# Completely disable OP_RETURN data carrier transactions
datacarrier=0
# Set OP_RETURN size to 0 (redundant with datacarrier=0 but ensures strictness)
datacarriersize=0
# Reject non-standard transactions (helps filter unusual spam patterns)
acceptnonstdtxn=0
# Lower mempool expiry time to clear spam faster (24 hours instead of default 336)
mempoolexpiry=24
# Reduce maximum transaction size to limit large data embedding
maxtxfee=0.1
# Limit ancestor/descendant chains (reduces complex spam patterns)
limitancestorcount=5
limitdescendantcount=5
limitancestorsize=50
limitdescendantsize=50

permitbaremultisig=0

# Create an RPC user/password
rpcuser=nodeadmin
rpcpassword=${PASSWORD}

# Enable ZMQ for Lightning compatibility
zmqpubrawblock=tcp://127.0.0.1:28332
zmqpubrawtx=tcp://127.0.0.1:28333

# Tor configuration
proxy=127.0.0.1:9050
listen=1
bind=127.0.0.1:8333

# Enable connections via Tor
listenonion=1

# Optional: Create a hidden service to make your node accessible
torcontrol=127.0.0.1:9051

# Force all connections through Tor
onlynet=onion

# Connection limits to reduce resource usage
maxconnections=40
maxuploadtarget=1000

# Ban nodes that send spam transactions
banscore=10
bantime=86400
```

# Basic Settings (explain)

### server=1: 
Enables the RPC server, allowing other applications (like Lightning nodes) to communicate with your Bitcoin node. This is essential for using the node with Lightning Network implementations.
### daemon=1:
Runs Bitcoin Core as a background service (daemon), so it doesn't require an active terminal session to keep running. It will continue running in the background.
### txindex=1: 
Creates and maintains a full transaction index. This makes it possible to query any transaction by its ID, not just those related to your wallet. Required for Lightning and other advanced functions.

## Data Directory

### datadir=/mnt/bitcoin: 
Specifies where Bitcoin Core will store all blockchain data, wallet files, and other data. Points to your external Samsung SSD mounted at /mnt/bitcoin.

## Network Settings

### listen=1: 
Allows your node to accept incoming connections from other Bitcoin nodes. This helps strengthen the network and may improve your connection quality.

## Performance Settings

### dbcache=4000: 
Allocates 4GB of RAM for the database cache. This significantly speeds up blockchain validation, especially during initial sync. With 16GB RAM, this is a reasonable allocation.
### maxmempool=500: 
Sets the maximum memory pool size to 500MB. The mempool stores unconfirmed transactions. Larger values use more RAM but can track more pending transactions.

## Initial Block Download Optimization

### blocksonly=1: 
During initial sync, only downloads blocks and not unconfirmed transactions. This significantly reduces bandwidth usage and speeds up the initial sync. You can set this to 0 after syncing is complete.

## RPC Authentication

### rpcuser=nodeadmin: 
Username for RPC (Remote Procedure Call) authentication. Used when other applications need to communicate with your Bitcoin node.
### rpcpassword=
YourVeryStrongPasswordHere: Password for RPC authentication. Critical for security - use a strong, unique password here.

## ZMQ Settings (for Lightning)

## zmqpubrawblock=tcp://127.0.0.1:28332:
 Enables ZeroMQ notifications for new blocks on port 28332. Lightning implementations use this to get immediate notifications when new blocks are found.
## zmqpubrawtx=tcp://127.0.0.1:28333: 
Enables ZeroMQ notifications for new transactions on port 28333. Lightning implementations use this to monitor the blockchain for relevant transactions.

## Tor

```shell
sudo nano /etc/tor/torrc

ControlPort 9051
CookieAuthentication 1
CookieAuthFile /var/run/tor/control.authcookie
HashedControlPassword 16:HASHVALUE  # Only if you set torpassword above

tor --hash-password "yourpassword"

# Check if your node is using Tor by running:
bitcoin-cli getnetworkinfo | grep torsm
```