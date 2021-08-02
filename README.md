# OpenData health

## sysctl configuration

Therefore, quic-go tries to increase the buffer size. The way to do this is an OS-specific, and we currently have an implementation for linux, windows and darwin. However, an application is only allowed to do increase the buffer size up to a maximum value set in the kernel. Unfortunately, on Linux this value is rather small, too small for high-bandwidth QUIC transfers.

```
sysctl -w net.core.rmem_max=2500000
```


## Run with docker
```sh
docker run -d \
  --name opendata-ipfs-node \
  --restart=always \
  -v /var/data/ipfs:/data/ipfs \
  -p 4001:4001 \
  -p 4001:4001/udp \
  --env OPENDATA_NODE=mynode \
  sethealth/opendata-ipfs-node:latest
```

## Run in Docker compose

```yaml
version: "3.8"
services:
  opendata_ipfs_node:
    image: sethealth/opendata-ipfs-node
    environment:
      OPENDATA_NODE: mynode

    ports:
      - "4001:4001" # ipfs swarm
      - "4001:4001/udp"

    volumes:
      - ./ipfs/ipfs0:/data/ipfs

```
