BINARY := stmt2redis
GO     := go

.PHONY: all build test lint clean

all: build

## build: compile the binary
build:
	$(GO) build -o $(BINARY) .

## test: run all unit tests
test:
	$(GO) test ./...

## lint: run go vet
lint:
	$(GO) vet ./...

## clean: remove compiled binary
clean:
	rm -f $(BINARY)
