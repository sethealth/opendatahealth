# OpenData health


## Run with docker
```
	docker run --name opendata-ipfs-node -d --restart=unless-stopped -v /var/data/ipfs:/data/ipfs -p 4001:4001 -p 4001:4001/udp sethealth/opendata-ipfs-node:latest
```

## Run in Docker compose

```
version: "3.8"
services:
  opendata_ipfs_node:
    image: sethealth/opendata-ipfs-node
    ports:
      - "4001:4001" # ipfs swarm
      - "4001:4001/udp"

    volumes:
      - ./ipfs/ipfs0:/data/ipfs
```
