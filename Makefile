recover-identity:
	echo $CLUSTER_IDENTITY > ipfs/cluster0/identity.json

deploy: recover-identity

run:
	docker-compose build --parallel

	docker-compose -f docker-compose.yaml up --renew-anon-volumes --remove-orphans
