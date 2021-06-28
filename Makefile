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

.PHONY: metalint
metalint:
ifeq (, $(shell command -v golangci-lint 2> /dev/null))
	@echo "===> Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.40.1
endif
	@echo "===> Running golangci-lint..."
	$$(go env GOPATH)/bin/golangci-lint run -c ./.golangci.yml ./...

.PHONY: generate
generate:
	@echo "==> Running generate..."
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/banking/banking.proto
