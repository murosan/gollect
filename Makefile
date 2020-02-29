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

format:
	go fmt .

lint: format
	go vet .
	staticcheck .
	golint .

# staticcheck fails on Github Actions
# https://github.com/murosan/gollect/runs/476703219?check_suite_focus=true
lint-ci:
	go vet .
	golint .
