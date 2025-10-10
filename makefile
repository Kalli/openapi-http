.PHONY: build run test clean

build:
	go build -o bin/openapi-http cmd/openapi-http/main.go

run:
	go run cmd/openapi-http/main.go

test:
	go test ./...

clean:
	rm -rf bin/

install:
	go install github.com/kalli/openapi-http

# example usage
example:
	go run cmd/openapi-http/main.go https://petstore3.swagger.io/api/v3/openapi.json

example-operation:
	go run cmd/openapi-http/main.go -operation-id addPet https://petstore3.swagger.io/api/v3/openapi.json 
