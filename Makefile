install:
	go install ./cmd/gollect.go

clean:
	go mod tidy

test:
	go test ./gollect

test-c:
	mkdir -p ./out
	go test -cover -coverprofile ./out/cover.out ./gollect
	go tool cover -html=./out/cover.out -o ./out/cover.html
	open ./out/cover.html
