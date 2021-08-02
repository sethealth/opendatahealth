build:
	docker build -t sethealth/opendata-ipfs-node .

run-tmp:
	docker run --rm --name opendata-ipfs-node \
		-v /tmp/data/ipfs:/data/ipfs \
		-p 4001:4001 \
		-p 5001:5001 \
		-p 4001:4001/udp \
		--env OPENDATA_NODE=mynode
		sethealth/opendata-ipfs-node:latest

run-prod:
	docker run -d \
		--name opendata-ipfs-node \
		--restart=always \
		-v /var/data/ipfs:/data/ipfs \
		-p 4001:4001 \
		-p 4001:4001/udp \
		--env OPENDATA_NODE=mynode \
		sethealth/opendata-ipfs-node:latest


run-prod:
	docker run -d --name opendata-ipfs-node --restart=always -v /var/data/ipfs:/data/ipfs -p 4001:4001 -p 4001:4001/udp --env OPENDATA_NODE=sethealth_do sethealth/opendata-ipfs-node:latest

run-compose: build
	docker-compose -f docker-compose.yaml up --renew-anon-volumes --remove-orphans

release: build
	docker image push sethealth/opendata-ipfs-node
	docker image tag sethealth/opendata-ipfs-node sethealth/opendata-ipfs-node:v0.0.3
