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

test-ci: clean
	go test -race .

lint:
	go fmt .
	go vet .
	staticcheck .
	golint .
