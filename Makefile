build:
	@go build -o bin/api
run: build
	@./bin/api
runOnPort:
	./bin/api --listenAddr :7000
seed:
	@go run scripts/seed.go
test:
	@go test -v ./...