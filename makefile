build:
	@go build -o bin/tojalB3 main.go

run: build
	@./bin/tojalB3

test:
	@go test -cover ./...
