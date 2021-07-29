build:
	docker build -t sethealth/opendata-ipfs-node .

run-tmp:
	docker run --rm --name opendata-ipfs-node \
		-v /tmp/data/ipfs:/data/ipfs \
		-p 4001:4001 \
		-p 5001:5001 \
		-p 4001:4001/udp \
		sethealth/opendata-ipfs-node:latest

run-prod:
	docker run --name opendata-ipfs-node -d --restart=unless-stopped -v /var/data/ipfs:/data/ipfs -p 4001:4001 -p 4001:4001/udp sethealth/opendata-ipfs-node:latest

run-compose:
	docker-compose build --parallel
	docker-compose -f docker-compose.yaml up --renew-anon-volumes --remove-orphans

release: build
	docker push sethealth/opendata-ipfs-node

