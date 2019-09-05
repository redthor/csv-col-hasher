APPLICATION_NAME=csv-col-hasher
GO_BIN=go
DOCKER_BIN=docker
VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')

.DEFAULT_GOAL := help
.PHONY: stop build build-alpine clean test help default

help:
	@echo 'Management commands for csv-col-hasher'
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

##
## Testing
##---------------------------------------------------------------------------

run:                     ## Run the program
	$(GO_BIN) run main.go -f tests/fixtures/test.csv -n 0

test:                    ## Run the tests
	$(GO_BIN) test ./...

##
## Project binary build
##---------------------------------------------------------------------------

get-deps:                ## Update the project's dependencies
	$(GO_BIN) get -u

build:                   ## Compile the binary
	@echo "building ${APPLICATION_NAME} ${VERSION}"
	$(GO_BIN) build -ldflags \
		"-X github.com/pepperstone/csv-col-hasher/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/pepperstone/csv-col-hasher/version.BuildDate=${BUILD_DATE}" \
		-o bin/${APPLICATION_NAME}

build-alpine:            ## Compile a binary optimised for Alpine
	@echo "building ${APPLICATION_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	$(GO_BIN) build -ldflags \
		'-w -linkmode external -extldflags "-static" -X github.com/pepperstone/csv-col-hasher/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/pepperstone/csv-col-hasher/version.BuildDate=${BUILD_DATE}' \
		-o bin/${APPLICATION_NAME}

clean:                   ## Remove the binary
	@test ! -e bin/${APPLICATION_NAME} || rm bin/${APPLICATION_NAME}
