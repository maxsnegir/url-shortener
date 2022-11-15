build:
	go build -v ./cmd/shortener/

run:
	./shortener

run-test:
	go test ./... -v -cover

docker-run:
	docker-compose up --build

docker-stop:
	docker-compose stop