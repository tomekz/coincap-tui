VERSION = $$(git describe --abbrev=0 --tags)
COMMIT_REV = $$(git rev-list -n 1 $(VERSION))

version:
	@echo $(VERSION)

commit_rev:
	@echo $(COMMIT_REV)

build:
	go build -ldflags "-X github.com/tomekz/coincap-tui/cmd.version=$(VERSION)" -o bin/tui

run: build
	./bin/tui

debug: build
	DEBUG=true ./bin/tui
