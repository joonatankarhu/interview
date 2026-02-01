.PHONY: test start stop run
test:
	go test ./... -v

# Start server in background (port 3000). PID saved to .server.pid.
start:
	go run ./cmd/... & echo $$! > .server.pid

# Stop server started with 'make start'.
stop:
	@kill $$(cat .server.pid 2>/dev/null) 2>/dev/null; rm -f .server.pid; echo "Server stopped"

# Run server in foreground (Ctrl+C to stop).
run:
	go run ./cmd/...
