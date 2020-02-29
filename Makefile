GOPATH = $(shell go env GOPATH)
GOBIN = $(GOPATH)/bin

install:
	go install ./cmd/gollect

install-tools:
	go mod download
	grep _ tools/tools.go | \
	awk '{print $$2}' | \
	xargs -tI % go install %

clean:
	go mod tidy
	go clean

test:
	go test .

test-c:
	mkdir -p ./out
	go test -cover -coverprofile ./out/cover.out .
	go tool cover -html=./out/cover.out -o ./out/cover.html
	open ./out/cover.html

lint: install-tools
	go fmt .
	go vet .
	$(GOBIN)/staticcheck .
	$(GOBIN)/golint .

ci: clean install-tools
	go vet .
	$(GOBIN)/staticcheck .
	$(GOBIN)/golint .
	go test -race .
