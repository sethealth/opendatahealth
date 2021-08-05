build:
	docker build -t sethealth/opendata-sentinel .

run-compose: build
	docker-compose -f docker-compose.yaml up

release: build
	docker image push sethealth/opendata-sentinel
	docker image tag sethealth/opendata-sentinel sethealth/opendata-sentinel:v0.0.1
