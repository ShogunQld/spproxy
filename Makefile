## help: Print this help message
.PHONY: help
help:
	@echo 'Sticky Port Proxy (spproxy) Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run: Run the Sticky Proxy Service from source code
.PHONY: run
run:
	go run cmd/main.go data/config.yaml

## test: Run all the tests
.PHONY: test
test:
	go test ./...

## build: Build proxy service executable for current architecture
.PHONY: build
build:
	go build -o bin/spproxy cmd/main.go

## build-linux: Build spproxy for linux amd64
.PHONY: build-linux
build-linux:
	env GOOS=linux GOARCH=amd64 go build -o bin/spproxy-linux-amd64 cmd/main.go

## build-mac: Build spproxy for MacOS arm64
.PHONY: build-mac
build-mac:
	env GOOS=darwin GOARCH=arm64 go build -o bin/spproxy-mac-arm64 cmd/main.go

## build-win: Build spproxy for Windows amd64
.PHONY: build-win
build-win:
	env GOOS=windows GOARCH=amd64 go build -o bin/spproxy.exe cmd/main.go

## build-all: Build spproxy for all architectures
.PHONY: build-all
build-all: build-linux build-win build-mac

## install: Install pre-requisite libraries: air
.PHONY: install
install:
	go install github.com/air-verse/air@latest

## start-test-spproxy: Start Test Sticky Port Proxy
.PHONY: start-test-spproxy
start-test-spproxy:
	@echo 'Start Test Sticky Port Proxy on 8080'
	@echo curl localhost:8080/ to hit the current sticky port
	@echo curl localhost:8080/server1/ to hit test server 1 and make port 8081 sticky
	@echo curl localhost:8080/server2/ to hit test server 2 and make port 8082 sticky
	@echo curl localhost:8080/server3/ to hit test server 3 and make port 8083 sticky
	go run cmd/main.go test/test-config.yaml

## start-test-server1: Start Test Server 1
.PHONY: start-test-server1
start-test-server1:
	@echo 'Start Test Server 1 on port 8081 - /server1 endpoint in test spproxy'
	go run test/main.go 8081

## start-test-server2: Start Test Server 2
.PHONY: start-test-server2
start-test-server2:
	@echo 'Start Test Server 2 on port 8082 - /server2 endpoint in test spproxy'
	go run test/main.go 8082

## start-test-server3: Start Test Server 3
.PHONY: start-test-server3
start-test-server3:
	@echo 'Start Test Server 3 on port 8083 - /server3 endpoint in test spproxy'
	go run test/main.go 8083

$(shell mkdir -p bin)
