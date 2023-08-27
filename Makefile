build:
	go build -o bin/tui

run: build
	./bin/tui

debug: build
	DEBUG=true ./bin/tui
