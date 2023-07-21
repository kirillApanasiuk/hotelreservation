build:
	@go build -o bin/api
run: build
	@./bin/api
runOnPort:
	./bin/api --listenAddr :7000
test:
	@go test -v ./...