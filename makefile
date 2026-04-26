build:
	@go build -o bin/tojalB3 cmd/main.go

run: build
	@./bin/tojalB3

test:
	@go test -cover ./...


