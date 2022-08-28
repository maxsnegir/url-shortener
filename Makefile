build:
	go build -v ./cmd/shortener/

run:
	./shortener

docker-run:
	docker-compose up --build

docker-stop:
	docker-compose stop