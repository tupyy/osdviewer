GIT_COMMIT=$(shell git rev-list -1 HEAD --abbrev-commit)

build:
	go build -mod=vendor $(BUILD_ARGS) -ldflags "-X main.CommitID=$(GIT_COMMIT) -s -w" \
	-o $(CURDIR)/bin/osdviewer $(CURDIR)/main.go
run:
	$(PWD)/bin/osdviewer

install:
	cp $(PWD)/bin/osdviewer $(HOME)/.local/bin/
