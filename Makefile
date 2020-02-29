install:
	go install ./cmd/gollect

install-tools:
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

format:
	go fmt .

format-keep:
	test -z "$(shell gofmt -s -l .| grep -Ev 'testdata/codes|out/')"

lint:
	go vet .
	staticcheck .
	golint .

lint-w: format lint

ci: clean install-tools format-keep lint
	go test -race .
