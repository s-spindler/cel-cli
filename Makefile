BINARY_NAME=cel-cli

build:
	mkdir -p bin
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows main.go

test:
	go test -v ./...

clean:
	go clean
	rm -f bin/*
