build:
	go build -o bin/tui

run: build
	./bin/tui
