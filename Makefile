GOPATH = $(shell go env GOPATH)
GOBIN = $(GOPATH)/bin

install:
	go install ./cmd/gollect

install-tools:
	go mod download
	# go install honnef.co/go/tools/cmd/staticcheck

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
	go fmt . ./testdata
	go vet .
	$(GOBIN)/staticcheck .

ci: clean install-tools
	go vet .
	# $(GOBIN)/staticcheck .
	go test -race .
