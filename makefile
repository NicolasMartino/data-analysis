test:
	go test -v -cover ./...
build:
	go build -o ./dist/ ./src/main.go 
.PHONY: test build