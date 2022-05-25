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

base: clean test swagger fmt build

build:
	go build -ldflags "$(LDFLAGS)" -mod vendor -o bin/weave main.go

run: 
	go run -mod vendor main.go

test:
	go test -ldflags -s -v --cover $(PKGS)

clean:
	@rm -rf bin/
	go mod tidy
	go mod vendor

fmt:
	gofmt -s -w $(GOFILES)

swagger:
	swag init

init: install-swagger postgres redis
	git config core.hooksPath .githooks
	echo "init all"

install-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest

postgres:
	@docker start mypostgres || docker run --name mypostgres -d -p 5432:5432 -e POSTGRES_PASSWORD=123456 postgres
	until docker exec mypostgres psql -U postgres; do echo "wait postgres start"; sleep 1; done
	cat scripts/db.sql | docker exec -i mypostgres psql -U postgres

exec-db:
	docker exec -it mypostgres psql -d weave -U postgres

redis:
	@docker start myredis || docker run --name myredis -d -p 6379:6379 redis --appendonly yes --requirepass 123456

ui:
	cd web && npm i && npm run dev
