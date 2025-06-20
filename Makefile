default: build

build:
	go build ./...

mod:
	go mod tidy

.PHONY: default build mod
