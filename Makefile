run:
	docker-compose build --parallel
	docker-compose -f docker-compose.yaml up --remove-orphans
