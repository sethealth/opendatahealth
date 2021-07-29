# OpenData health


## Run with docker
```sh
docker run
  -d \
  --name opendata-ipfs-node \
  --restart=always \
  -v /var/data/ipfs:/data/ipfs \
  -p 4001:4001 \
  -p 4001:4001/udp \
  -env OPENDATA_NODE=mynode \
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
