build:
	@go build -o bin/spinner

run: build
	@./bin/spinner

test:
	@go test -v ./...