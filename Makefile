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

clean:
	go fmt ./...
	go mod tidy -v
	go mod verify
	go vet ./...

test:
	go test ./... -coverprofile=cover.out -timeout 2s
	go tool cover --html=cover.out
