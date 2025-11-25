default: build

VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT  ?= $(shell git rev-parse HEAD)
DATE    ?= $(shell date -u "+%Y-%m-%dT%H:%M:%SZ")

LDFLAGS = -s -w \
	  -X main.buildVersion=$(VERSION) \
	  -X main.buildDate=$(DATE) \
	  -X main.buildCommit=$(COMMIT)

build:
	go build -o tenuki -ldflags="$(LDFLAGS)" .

mod:
	go get -u github.com/ymattw/googs
	go mod tidy

release:
	goreleaser release --clean

snapshot:
	goreleaser release --clean --snapshot

.PHONY: default build mod release snapshot
