base: clean swagger
	go build -o bin/weave main.go

run: 
	go run main.go

TESTPKG=$(shell go list ./...)
test:
	go test -ldflags -s -v --cover $(TESTPKG)

clean:
	@rm -rf bin/
	go mod tidy
	go mod vendor

swagger:
	swag init

init: install-swagger postgres redis
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
