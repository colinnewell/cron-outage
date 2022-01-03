# FIXME: it would be nice to encode branch too
VERSION  := $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)

all: cron-outage

cron-outage: cmd/cron-outage/*.go internal/cron/*
	go build -o cron-outage -ldflags "-X main.Version=$(VERSION)" cmd/cron-outage/*.go

test:
	go test ./...

install: cron-outage
	cp cron-outage /usr/local/bin

lint:
	golangci-lint run
	./ensure-gofmt.sh

license-check:
	# gem install license_finder
	license_finder
