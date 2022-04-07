.DEFAULT_GOAL := help
BIN:=$(CURDIR)/bin
BIN_LINTER:=$(BIN)/golangci-lint
VERSION:=$(shell cat VERSION)
REGISTRY_DOMAIN=ghcr.io
REGISTRY_NAME=ghcr.io/wault-pw/alice

help:
	@echo 'Available targets: $(VERSION)'
	@echo '  make build'
	@echo '  make push'
	@echo '  make outdated'
	@echo ' '
	@echo '  make db:status'
	@echo '  make db:up'
	@echo '  make db:down'
	@echo '  make NAME="create_users" db:create'
	@echo ' '
	@echo '  make test'
	@echo '  make outdated'

.PHONY: install-lint
install-lint:
ifeq ($(wildcard $(BIN_LINTER)),)
	$(info Downloading golangci-lint)
	GOBIN=$(BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.1
endif

t:
	$(eval PG_DSN?=postgres://localhost:5432/alice_test?sslmode=disable&timezone=utc)

d:
	$(eval PG_DSN?=postgres://localhost:5432/alice?sslmode=disable&timezone=utc)

db\:up: d
db\:up:
	goose -dir migrations postgres "$(PG_DSN)" up

db\:down: d
db\:down:
	goose -dir migrations postgres "$(PG_DSN)" down

db\:status: d
db\:status:
	goose -dir migrations postgres "$(PG_DSN)" status

db\:create: NAME=$NAME
db\:create:
	goose -dir migrations postgres "$(PG_DSN)" create $(NAME) sql

proto:
	protoc --proto_path=protos --go_out=. alice_v1.proto

test: t
test:
	PG_DSN="$(PG_DSN)" go test -count=1 -p 4 -race -cover -covermode atomic ./...

outdated:
	go list -u -m all

lint: install-lint
lint:
	$(BIN_LINTER) run --config=.golangci.yaml ./...

generate:
	go generate ./...

generate_build:
	go generate cmd/goose.go

linux: generate_build
linux:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -ldflags "-s -X 'main.Version=$(VERSION)'" -a -installsuffix cgo -o build/linux

mac: generate_build
mac:
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -ldflags "-s -X 'main.Version=$(VERSION)'" -a -installsuffix cgo -o build/mac

.PHONY: build
build: export TAG=$(VERSION)
build:
	docker build --no-cache -f ./Dockerfile -t ${REGISTRY_NAME}:${TAG} .
	docker tag ${REGISTRY_NAME}:${TAG} ${REGISTRY_NAME}:latest

push: export TAG=$(VERSION)
push:
	docker push ${REGISTRY_NAME}:${TAG}
	docker push ${REGISTRY_NAME}:latest
