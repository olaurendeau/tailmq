build: build-darwin-amd64 build-linux-amd64

install: 
	go install

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o build/tailmq-darwin-amd64

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o build/tailmq-linux-amd64
