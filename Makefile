CLIENT_MAIN=./cmd/client
SERVER_MAIN=./cmd/server
BINARY_CLIENT=grpc-client
BINARY_SERVER=grpc-server

.PHONY:build
build: ## Build GRPC client and server
	@echo Building client... && \
	go build -o $(BINARY_CLIENT) $(CLIENT_MAIN)
	@echo Building server... && \
	go build -o $(BINARY_SERVER) $(SERVER_MAIN)

.PHONY:clean
clean: ## Remove binaries and run go clean
	@echo Cleaning... && \
	rm -f $(BINARY_CLIENT) \
    rm -f $(BINARY_SERVER) \
	go clean

.PHONY:gen
gen: ## Generate code using the protocol buffer compiler
	go generate .

.PHONY: help
help: ## Show all commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
