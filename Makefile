build:
	docker build -t opendata-node .

run:
	docker run --rm --name opendata-node \
		-v /tmp/ipfs-opendata-staging:/export \
		-v /tmp/ipfs-opendata-data:/data/ipfs \
		-p 8080:8080 \
		-p 4001:4001 \
		-p 127.0.0.1:5001:5001 \
		opendata-node:latest

compose:
	docker-compose build --parallel
	docker-compose -f docker-compose.yaml up --renew-anon-volumes --remove-orphans
