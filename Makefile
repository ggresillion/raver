BACKEND=backend
FRONTEND=frontend

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all

all: help

## Build:
build: build-backend build-frontend ## Build the backend and the frontend

build-backend: ## Build the backend
	make -C ${BACKEND} build

build-frontend: ## Build the frontend
	make -C ${FRONTEND} build

clean: clean-backend clean-frontend ## Remove build related file

clean-backend: ## Remove backend build related file
	make -C ${BACKEND} clean

clean-frontend: ## Remove frontend build related file
	make -C ${FRONTEND} clean

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