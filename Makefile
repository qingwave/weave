.PHONY: base build run test clean swagger init ui

GIT_VERSION = $(shell describe --tags 2>/dev/null)
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GIT_STATE = $(shell [ -z $(git status --porcelain 2>/dev/null) ] && echo "dirty" || echo "clean")
BUILD_DATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS = -X github.com/qingwave/weave/pkg/version.gitVersion=$(GIT_VERSION) \
	-X github.com/qingwave/weave/pkg/version.gitCommit=$(GIT_COMMIT) \
	-X github.com/qingwave/weave/pkg/version.gitBranch=$(GIT_BRANCH) \
	-X github.com/qingwave/weave/pkg/version.gitTreeState=$(GIT_STATE) \
	-X github.com/qingwave/weave/pkg/version.buildDate=$(BUILD_DATE)

PKGS = $(shell go list ./...)
GOFILES = $(shell find . -name "*.go" -type f -not -path "./vendor/*")

base: clean test lint swagger fmt build

help: ## display the help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: ## build server
	go build -ldflags "$(LDFLAGS)" -mod vendor -o bin/weave main.go

run: ## run server
	go run -mod vendor main.go

test: ## run unit test
	go test -ldflags -s -v -coverprofile=cover.out $(PKGS)
	go tool cover -func=cover.out -o coverage.txt

lint: install-golangci-lint ## run golangci lint
	golangci-lint run

clean: ## clean bin and go mod
	@rm -rf bin/
	go mod tidy
	go mod vendor

fmt: ## golang format
	gofmt -s -w $(GOFILES)

swagger: ## swagger init
	swag init

init: install-swagger install-golangci-lint postgres redis ## install all dependencies and init config
	git config core.hooksPath .githooks
	echo "init all"

install-swagger: ## install swagger from golang
	go install github.com/swaggo/swag/cmd/swag@v1.8.12

install-golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

postgres: ## init postgres db
	@docker start mypostgres || docker run --name mypostgres -d -p 5432:5432 -e POSTGRES_PASSWORD=123456 postgres
	until docker exec mypostgres psql -U postgres; do echo "wait postgres start"; sleep 1; done
	cat scripts/db.sql | docker exec -i mypostgres psql -U postgres

exec-db: ## exec to db container
	docker exec -it mypostgres psql -d weave -U postgres

redis: ## init redis
	@docker start myredis || docker run --name myredis -d -p 6379:6379 redis --appendonly yes --requirepass 123456

ui: ## run ui locally
	cd web && npm i && npm run dev

SERVER_IMG=qingwave/weave-server
docker-build-server: ## build server image
	docker build -t $(SERVER_IMG) .

docker-run-server: ## run server in docker
	docker run --network host -v $(shell pwd)/config:/config -v /var/run/docker.sock:/var/run/docker.sock $(SERVER_IMG)

FRONENT_IMG=qingwave/weave-frontend
docker-build-ui: ## build frontend image
	cd web && docker build --build-arg BUILD_OPTS=build -t $(FRONENT_IMG) .

MOCK_FRONENT_IMG=qingwave/weave-frontend:mock
docker-build-ui-mock: ## build mock frontend image
	cd web && docker build --build-arg BUILD_OPTS=build-with-mock -t $(MOCK_FRONENT_IMG) .

docker-run-ui-mock: ## run mock frontend in docker
	docker run --network host -v $(shell pwd)/web/nginx.conf:/etc/nginx/conf.d/nginx.conf $(MOCK_FRONENT_IMG)
