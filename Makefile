.PHONY: all
all: build

.PHONY: build build_server build_client
build: build_server build_client

build_server:
	mkdir -p bin
	go build -o bin/server ./cmd/server/server.go

build_client:
	mkdir -p bin
	go build -o bin/client ./cmd/client/client.go

.PHONY: test test_server test_client
test: test_server test_client

test_server:
	go test -v -run '^Test' ./cmd/server

test_client:
	go test -v -run '^Test' ./cmd/client

.PHONY: clean
clean:
	rm -f bin/server bin/client
	go clean