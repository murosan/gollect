install:
	go install cmd

clean:
	go mod tidy

test-c:
	mkdir -p ./out
	go test -cover -coverprofile ./out/cover.out ./gollect
	go tool cover -html=./out/cover.out -o ./out/cover.html
	open ./out/cover.html
