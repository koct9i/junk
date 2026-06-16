all: fmt vet test build

build:
	go build .

install:
	go install .

lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

generate:
	go generate ./...
