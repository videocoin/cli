.PHONY: all clean

export CGO_ENABLED=1
export GOOS=darwin
export GOARCH=amd64

VERSION=0.1.1
BUILD=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"
OUTDIR=build

all: cli mrtmp

cli:
	go build -o $(OUTDIR)/cli $(LDFLAGS) ./cmd/cli 

mrtmp:
	go build -o $(OUTDIR)/minirtmp ./cmd/minirtmp

clean:
	-rm -r $(OUTDIR)