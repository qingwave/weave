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

install-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest

postgres:
	docker run --name mypostgres -d -p 5432:5432 -e POSTGRES_PASSWORD=123456 postgres
	
exec-db:
	docker exec -it mypostgres psql -d weave -U postgres
