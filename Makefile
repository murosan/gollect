install:
	go install ./cmd/gollect

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

lint:
	go fmt .
	go vet .
	staticcheck .
	golint .

ci:
	go mod download
	go vet .
	golint -set_exit_status .
	staticcheck .
	test -z "$(shell gofmt -s -l .| grep -Ev 'testdata/codes|out/')"
	go test -race .
