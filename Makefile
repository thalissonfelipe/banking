NAME=banking
OS ?= linux

.PHONY: dev-docker
dev-docker:
	@echo "===> Starting server application..."
	docker-compose up --build

.PHONY: dev-local
dev-local:
	@echo "===> Starting server application..."
	go mod download
	docker-compose up --d db
	go run cmd/api/main.go

.PHONY: test
test:
	@echo "==> Running tests..."
	go test -v ./...
