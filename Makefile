GO=go
GOFLAGS=
BIN=triageprof
GO_PLUGIN=go-pprof-http
PYTHON_PLUGIN=python-cprofile
NODE_PLUGIN=node-inspector

.PHONY: all build test demo demo-python clean

all: build

build:
	$(GO) build $(GOFLAGS) -o bin/$(BIN) ./cmd/triageprof
	$(GO) build $(GOFLAGS) -o plugins/bin/$(GO_PLUGIN) ./plugins/src/$(GO_PLUGIN)
	$(GO) build $(GOFLAGS) -o plugins/bin/$(NODE_PLUGIN) ./plugins/src/$(NODE_PLUGIN)
	chmod +x plugins/src/$(PYTHON_PLUGIN)/main.py
	cp plugins/src/$(PYTHON_PLUGIN)/main.py plugins/bin/$(PYTHON_PLUGIN)

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
	bin/$(BIN) run --plugin $(GO_PLUGIN) --target-url http://localhost:6060 --duration 5 --outdir out
	
	# Cleanup
	kill $$SERVER_PID || true
	
	echo "Demo completed. Results in out/ directory."

demo-python: build
	# Start Python demo server in background
	cd examples/python-demo-server && python3 demo.py &
	SERVER_PID=$$!
	echo "Python demo server started on PID $$SERVER_PID"
	
	# Wait for server to start
	sleep 2
	
	# Generate load on Python server
	curl -s http://localhost:8080/cpu-hotspot > /dev/null &
	curl -s http://localhost:8080/alloc-heavy > /dev/null &
	curl -s http://localhost:8080/memory-leak > /dev/null &
	wait
	
	# Run triageprof with Python plugin
	mkdir -p out-python
	bin/$(BIN) run --plugin $(PYTHON_PLUGIN) --target-type python --target-command "python3 ../../examples/python-demo-server/demo.py" --duration 5 --outdir out-python
	
	# Cleanup
	kill $$SERVER_PID || true
	
	echo "Python demo completed. Results in out-python/ directory."

demo-node: build
	# Start Node.js demo server in background
	cd examples/node-demo-server && node server.js &
	SERVER_PID=$$!
	echo "Node.js demo server started on PID $$SERVER_PID"
	
	# Wait for server to start
	sleep 2
	
	# Generate load on Node.js server
	curl -s http://localhost:3000 > /dev/null &
	curl -s http://localhost:3000 > /dev/null &
	curl -s http://localhost:3000 > /dev/null &
	wait
	
	# Run triageprof with Node.js plugin
	mkdir -p out-node
	bin/$(BIN) run --plugin $(NODE_PLUGIN) --target-type node --target-command "node ../../examples/node-demo-server/server.js" --duration 5 --outdir out-node
	
	# Cleanup
	kill $$SERVER_PID || true
	
	echo "Node.js demo completed. Results in out-node/ directory."

clean:
	rm -rf bin/ plugins/bin/ out/

install:
	mkdir -p bin/ plugins/bin/