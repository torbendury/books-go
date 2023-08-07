CGO_ENABLED:=1

startpsql:
	cd hack
	docker compose -f hack/docker-compose.yml down
	docker volume prune -f
	docker compose -f hack/docker-compose.yml up -d
	cd ..

stoppsql:
	cd hack
	docker compose -f hack/docker-compose.yml down
	docker volume prune -f
	cd ..

run:
	go run cmd/main.go

runpg:
	go run cmd/main.go -postgres

compose:
	docker compose -f hack/docker-compose.yml down
	docker volume prune -f
	docker image prune -f
	docker compose -f hack/docker-compose.yml up --build

clean:
	go fmt ./...
	go mod tidy -v
	go mod verify
	go vet ./...
	go clean

puml:
	go-plantuml generate -d . -r

test:
	go test -race ./... -coverprofile=cover.out -timeout 2s
	go tool cover --html=cover.out

build:
	docker build -t torbendury/books-go:test --target test .
	docker build -t torbendury/books-go:latest .

all: clean test build
