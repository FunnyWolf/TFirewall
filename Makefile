CLIENT_SOURCE=./cmd/client
SERVER_SOURCE=./cmd/server

LDFLAGS="-s -w"
GCFLAGS="all=-trimpath=$GOPATH"

CLIENT_BINARY=tfc
SERVER_BINARY=tfs
TAGS=release

OSARCH = "linux/amd64 linux/386 windows/amd64 windows/386"


.DEFAULT: help

dep: ## Install dependencies
	go get -d -v ./...
	go get -u github.com/mitchellh/gox


build: ## Build for the current architecture.
	go build -ldflags $(LDFLAGS) -gcflags $(GCFLAGS) -tags $(TAGS) -o bin/$(CLIENT_BINARY) $(CLIENT_SOURCE) && \
	go build -ldflags $(LDFLAGS) -gcflags $(GCFLAGS) -tags $(TAGS) -o bin/$(SERVER_BINARY) $(SERVER_SOURCE)

build-all: ## Build for every architectures.
	gox -osarch=$(OSARCH) -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -tags $(TAGS) -output "bin/$(SERVER_BINARY)_{{.OS}}_{{.Arch}}" $(SERVER_SOURCE)
	gox -osarch=$(OSARCH) -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -tags $(TAGS) -output "bin/$(CLIENT_BINARY)_{{.OS}}_{{.Arch}}" $(CLIENT_SOURCE)

clean:
	rm -rf certs
	rm bin/$(SERVER_BINARY)_*
	rm bin/$(CLIENT_BINARY)_*
