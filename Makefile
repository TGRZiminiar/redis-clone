run: build
	@./bin/redis-clone --listenAddr :5001
build:
	@go build -o bin/redis-clone

tp:
	go test proto_test.go -v


# go run main.go