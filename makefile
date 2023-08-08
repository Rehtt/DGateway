VERSION=$(shell git describe --tags --always)

.PHONY: build
build:
	go mod tidy && go build -o ./bin/run -tags=jsoniter -ldflags "-X main.Version=$(VERSION)" .