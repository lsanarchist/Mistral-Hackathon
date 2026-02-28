GO=go
GOFLAGS=
BIN=triageprof
PLUGIN=go-pprof-http

.PHONY: all build test demo clean

all: build

build:
	$(GO) build $(GOFLAGS) -o bin/$(BIN) ./cmd/triageprof
	$(GO) build $(GOFLAGS) -o plugins/bin/$(PLUGIN) ./plugins/src/$(PLUGIN)

test:
	$(GO) test $(GOFLAGS) ./...

demo: build
	# Start demo server in background
	cd examples/demo-server && $(GO) run main.go &
	SERVER_PID=$$!
	echo "Demo server started on PID $$SERVER_PID"
	
	# Wait for server to start
	sleep 2
	
	# Generate load
	./examples/load.sh
	
	# Run triageprof
	mkdir -p out
	bin/$(BIN) run --plugin $(PLUGIN) --target-url http://localhost:6060 --duration 5 --outdir out
	
	# Cleanup
	kill $$SERVER_PID || true
	
	echo "Demo completed. Results in out/ directory."

clean:
	rm -rf bin/ plugins/bin/ out/

install:
	mkdir -p bin/ plugins/bin/