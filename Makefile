.PHONY:  build clean dev test release

VERSION=0.1.0

default: build

build: clean
	CGO_ENABLED=0 go build -o ./dist/mikrotik-ros-exporter -a -ldflags '-s' -installsuffix cgo cmd/mikrotik-ros-exporter/main.go

clean:
	rm -rf ./dist/*

dev:
	gow run cmd/mikrotik-ros-exporter/main.go

test:
	go test -v ./...

release: clean
	goreleaser release
