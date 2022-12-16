build:
	go build -v ./cmd/shortener/

run:
	./shortener

gen-shortener-mock:
	mockgen -destination=internal/mocks/mock_storage.go -package=mocks github.com/maxsnegir/url-shortener/internal/storage ShortenerStorage

run-test:
	go test ./... -v -cover -race

docker-run:
	docker-compose up --build

docker-stop:
	docker-compose stop