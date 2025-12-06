GOPATH = $(shell go env GOPATH)
GOBIN = $(GOPATH)/bin
.PHONY: install install-tools clean test test-c ci

install:
	go install ./cmd/gollect

install-tools:
	go mod download

clean:
	go mod tidy
	go clean

test:
	go test -count=1 .

test-c:
	mkdir -p ./out
	go test -cover -coverprofile ./out/cover.out .
	go tool cover -html=./out/cover.out -o ./out/cover.html

ci: clean install-tools
	go test -race .
