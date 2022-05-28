.DELETE_ON_ERROR: clean

EXECUTABLES = go
K := $(foreach exec,$(EXECUTABLES),\
  $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

PROJECT_DEPENDENCIES := $(shell go list -m -f '{{if not (or .Indirect .Main)}}{{.Path}}{{end}}' all)

# avoid mocks in tests
GO_FILES       := $(shell go list ./...)

all: clean test

mod-update: tidy
	$(foreach dep, $(PROJECT_DEPENDENCIES), $(shell go get -u $(dep)))
	go mod tidy

tidy:
	go mod tidy

fmt:
	@go fmt $(GO_FILES)

vet:
	go vet $(GO_FILES)

lint:
	golangci-lint run

generate:
	go generate $(GO_FILES)

test: tidy fmt vet
	go test -race -covermode=atomic -coverprofile coverage.out -tags=unit $(GO_FILES)

test-coverage: test
	go tool cover -html=coverage.out

bench:
	go test -bench=. -benchmem

clean:
	rm -rf ./*.out
