# OpenData health

## sysctl configuration

Therefore, quic-go tries to increase the buffer size. The way to do this is an OS-specific, and we currently have an implementation for linux, windows and darwin. However, an application is only allowed to do increase the buffer size up to a maximum value set in the kernel. Unfortunately, on Linux this value is rather small, too small for high-bandwidth QUIC transfers.

```
sysctl -w net.core.rmem_max=2500000
```


## Run in Docker compose

```yaml
version: "3.8"
services:
  opendata_ipfs:
    image: ipfs/go-ipfs:v0.9.1
    restart: always
    sysctls:
      net.core.rmem_max: 2500000

    environment:
      IPFS_PATH: /data/ipfs
      IPFS_PROFILE: server,flatfs

    ports:
      - "4001:4001" # ipfs swarm
      - "4001:4001/udp"

    volumes:
      - ./ipfs/ipfs0:/data/ipfs

  opendata_sentinel:
    image: sethealth/opendata-sentinel
    restart: always
    environment:
      IPFS_URL: opendata_ipfs:5001
      PINSET_URL: https://api.set.health/pinset/openview-health/bu-5643105772503040
      OPENDATA_NODE: sethealth_do
    depends_on:
      - opendata_ipfs
```
