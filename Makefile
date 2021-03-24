BIN_NAME := $(or $(PROJECT_NAME),'libreria')
PKG_PATH := $(or $(PKG),'.')
PKG_LIST := $(shell go list ${PKG_PATH}/... | grep -v /vendor/)
GOLINT := golangci-lint

MIGRATE=migrate -path storage/postgres/migrations -database postgres://postgres:12345@localhost:5432/libreria?sslmode=disable

.PHONY: all
all: run

run: build-docker migrate-up

build: dep ## Build the binary file
	go build -o ./bin/${BIN_NAME} -a .

build-docker:
	docker-compose down
	docker-compose up -d --build
	./scripts/wait_then_run.sh localhost 5432 30

.PHONY: dep
dep: # Download required dependencies
	go mod tidy
	go mod download
	go mod vendor

## make migrate-create NAME="table"
migrate-create: ## Create migration file with name
	migrate create -ext sql -dir storage/postgres/migrations -format 20060102150405 $(NAME)

migrate-up: ## Run migrations
	$(MIGRATE) up

migrate-down: ## Rollback migrations
	$(MIGRATE) down

test-integration: dep ## Run integration tests
	docker-compose -f docker-compose-test.yml down
	docker-compose -f docker-compose-test.yml up -d
	POSTGRES_TEST_USER=postgres POSTGRES_TEST_NAME=libreria_test go test -v -tags integration -race -count=1 ./...
	docker-compose -f docker-compose-test.yml down

test-unit: dep ## Run unit tests
	go test -tags=unit  -race -count=1 -short ./...

clean: ## Remove previous build
	rm -f bin/$(BIN_NAME)

.PHONY: gen
gen:
	mockgen -package mock -source service/auth/service.go -destination service/auth/mock/service.go
	mockgen -package mock -source service/scheduler/service.go -destination service/scheduler/mock/service.go
	mockgen -package mock -source service/service.go -destination service/mock/service.go
	mockgen -package mock -source service/user/service.go -destination service/user/mock/service.go

check-lint:
	@which $(GOLINT) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.25.0

.PHONY: generate-mocks
generate-mocks: check-mockgen
	mockgen -package mock -source server/http/handlers/book.go -destination server/http/handlers/mock/book.go

.PHONY: check-mockgen
check-mockgen:
	@which mockgen || go get github.com/golang/mock/mockgen@v1.4.4
