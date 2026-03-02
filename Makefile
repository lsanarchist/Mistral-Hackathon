GO=go
GOFLAGS=
BIN=triageprof
GO_PLUGIN=go-pprof-http

.PHONY: all build test clean lint release

all: build

build:
	$(GO) build $(GOFLAGS) -o bin/$(BIN) ./cmd/triageprof
	$(GO) build $(GOFLAGS) -o plugins/bin/$(GO_PLUGIN) ./plugins/src/$(GO_PLUGIN)

test:
	$(GO) test $(GOFLAGS) ./...

# Run linter
lint:
	gofmt -l .
	golangci-lint run

# Create release artifacts
release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o release/triageprof-linux-amd64 ./cmd/triageprof
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o release/triageprof-darwin-amd64 ./cmd/triageprof
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o release/triageprof-windows-amd64.exe ./cmd/triageprof
	@echo "Release artifacts created in release/"

demo: build
	# Start demo server and run triageprof against it
	examples/demo-server/main &
	sleep 1
	mkdir -p demo-output
	bin/$(BIN) run --plugin $(GO_PLUGIN) --target-url http://localhost:6060 --duration 10 --outdir demo-output

clean:
	rm -rf bin/ plugins/bin/ out/
	@echo "Cleaned build artifacts"

help:
	@echo "Makefile targets:"
	@echo "  build   - Build binary and go-pprof-http plugin"
	@echo "  test    - Run all tests"
	@echo "  demo    - Run demo against built-in demo server"
	@echo "  clean   - Remove build artifacts"
	@echo "  release - Cross-compile release binaries"
	@echo "  help    - Show this help"

install:
	mkdir -p bin/ plugins/bin/