GO_CMD=go
BINARY_NAME=discordsoundboard
OUT_DIR=./out
BIN_DIR=${OUT_DIR}/bin
MAIN=cmd/main.go

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all

all: help

## Build:
build: ## Build your project and put the output binary in out/bin/
	mkdir -p ${BIN_DIR}
	$(GO_CMD) build -o ${BIN_DIR}/$(BINARY_NAME) ${MAIN}

clean: ## Remove build related file
	rm -fr ${OUT_DIR}

## Doc:
doc: ## Generates the doc from the openapi spec
	swag init -o ./api/docs -g api/api.go

## Format:
format: format-code format-doc ## Formats the code
format-code: ## Formats the go code
	go fmt ./...
format-doc: ## Formats the openapi docs
	swag fmt -g api/api.go

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)