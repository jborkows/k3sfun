.PHONY: build run dev test fmt migrate-up migrate-down clean gen sqlc templ sqlc-check docker

APP_NAME := shopping
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)

BUILD_VERSION ?= $(shell date +%s)
DOCKER_IMAGE ?= jborkows/shoppinglist

GOCACHE ?= $(CURDIR)/tmp/gocache
GOMODCACHE ?= $(CURDIR)/tmp/go-mod

export GOCACHE
export GOMODCACHE

SQLC ?= sqlc
TEMPL_VERSION ?= v0.3.960

SQLC_OUT := internal/db/db.go internal/db/models.go internal/db/queries.sql.go
SQLC_IN := internal/db/queries.sql sqlc.yaml

TEMPL_SRCS := $(shell find ./internal/web/views -name '*.templ')
TEMPL_STAMP := tmp/.templ.stamp

gen: sqlc templ

sqlc-check:
	@command -v $(SQLC) >/dev/null || (echo "sqlc not found; install it (e.g. go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)"; exit 1)

$(SQLC_OUT): $(SQLC_IN)
	$(MAKE) sqlc-check
	$(SQLC) generate

sqlc: $(SQLC_OUT)

$(TEMPL_STAMP): $(TEMPL_SRCS)
	go run github.com/a-h/templ/cmd/templ@$(TEMPL_VERSION) generate -path ./internal/web/views
	mkdir -p $(@D)
	touch $@

templ:
ifneq ($(strip $(TEMPL_SRCS)),)
	$(MAKE) $(TEMPL_STAMP)
else
	@true
endif

build: gen
	mkdir -p $(BIN_DIR) $(GOCACHE)
	go build -ldflags "-X main.buildVersion=$(BUILD_VERSION)" -o $(BIN) ./cmd/shopping

docker:
	docker build \
		--build-arg BUILD_VERSION=$(BUILD_VERSION) \
		-t $(DOCKER_IMAGE):$(BUILD_VERSION) \
		-t $(DOCKER_IMAGE):latest \
		.

run: gen
	mkdir -p $(GOCACHE)
	go run ./cmd/shopping

dev:
	air

test: gen
	mkdir -p $(GOCACHE)
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

clean:
	rm -rf ./tmp ./bin

migrate-up:
	migrate -path ./migrations -database "sqlite3://data/shopping.db" up

migrate-down:
	migrate -path ./migrations -database "sqlite3://data/shopping.db" down 1
