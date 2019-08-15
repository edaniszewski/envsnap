#
# envsnap - development targets
#

BIN_NAME    := envsnap
BIN_VERSION := 0.1.0
IMG_NAME	:= edaniszewski/envsnap

.PHONY: build clean cover docker fmt github-tag lint test version help


build:  ## Build the binary
	go build cmd/envsnap.go

clean:  ## Clean build and test artifacts
	@rm envsnap
	@rm coverage.out

cover:  ## Open a coverage report
	go tool cover -html=coverage.out

docker:  ## Build the docker image
	docker build -t ${IMG_NAME} .

fmt:  ## Run goimports formatting on all go files
	@find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

github-tag:  ## Create and push a tag with the current plugin version
	git tag -a ${BIN_VERSION} -m "${BIN_NAME} version ${BIN_VERSION}"
	git push -u origin ${BIN_VERSION}

lint:  ## Lint project source files
	@golint -set_exit_status ./cmd/...
	@golint -set_exit_status ./pkg/...

test:  ## Run project unit tests
	go test --race -coverprofile=coverage.out -covermode=atomic ./...

version:  ## Print the version of the project
	@echo "${BIN_VERSION}"

help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help
