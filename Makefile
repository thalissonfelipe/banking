.PHONY: dev-docker
dev-docker:
	@echo "===> Starting application on docker..."
	docker-compose up --build

.PHONY: dev-local
dev-local:
	@echo "===> Starting application locally..."
	go mod download
	docker-compose up --d db
	go run cmd/api/main.go

.PHONY: test
test:
	@echo "==> Running tests..."
	go test -v -race ./...

.PHONY: lint
lint:
ifeq (, $(shell command -v golangci-lint 2> /dev/null))
	@echo "===> Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.46.2
endif
	@echo "===> Running golangci-lint..."
	$$(go env GOPATH)/bin/golangci-lint run -c ./.golangci.yml ./...

.PHONY: generate
generate:
	@echo "==> Running generate..."
	swag init -g router.go -d banking/gateway/http -o docs/swagger
