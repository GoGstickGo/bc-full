GO ?= go
GOPATH := $(shell $(GO) env GOPATH)

all: build

build:
	$(GO) build -o node-dashboard ./cmd/server

run:
	$(GO) run main.go

clean:
	rm  node-dashboard

test:
	$(GO) test -v ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

autossh:
 # autossh -M 20000 -L ${RPCIP}:localhost:${RPCIP} -p ${SSHPORT} ${USER}@${HOST} -N -f

connect:
#	 ./node-dashboard -rpcuser=${USER}-rpcpass=${PW} -rpchost=localhost:${RPCIP} -listen=localhost:${HTTPSERVERIP}

.PHONY: all build run clean test fmt vet connect autossh