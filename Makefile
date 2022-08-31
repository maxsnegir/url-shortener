build:
	go build -v ./cmd/shortener/

run:
	./shortener

test:
	go test ./... -v

docker-run:
	docker-compose up --build

docker-stop:
	docker-compose stop

# Для локальной разработки поднимаем только базу
docker-run-dev:
	docker-compose up redis