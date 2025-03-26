```shell
wget https://bitcoincore.org/bin/bitcoin-core-28.1/SHA256SUMS

sha256sum --check SHA256SUMS --ignore-missing

tar -xzf bitcoin-25.1-x86_64-linux-gnu.tar.gz

sudo install -m 0755 -o root -g root -t /usr/local/bin bitcoin-25.1/bin/*

bitcoind --version
```
