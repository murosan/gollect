install:
	go install ./cmd/gollect

install-tools:
	grep _ tools/tools.go | \
	awk '{print $$2}' | \
	xargs -tI % sh -c 'go install % && go get -u %'

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

test-ci: clean
	go test -race .

format:
	go fmt .

format-keep:
	test -z "$(shell gofmt -s -l .| grep -Ev 'testdata/codes|out/')"

lint-ci:
	go vet .
	staticcheck .
	golint .

lint: format lint-ci
