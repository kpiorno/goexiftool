GOPATH:=$(shell go env GOPATH)

.PHONY: build
build: proto
	go build -o srv *.go

.PHONY: test
test:
	go test -v ./... -cover


local: test
	go mod tidy
	go run main.go

vet:
	go vet -v ./...

fmt:
	gofmt -w .

