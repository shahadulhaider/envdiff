BINARY := envdiff
PKG    := ./cmd/envdiff
VERSION ?= dev
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build install clean vet lint cross

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(PKG)

install:
	go install -ldflags "$(LDFLAGS)" $(PKG)

clean:
	rm -f $(BINARY)
	rm -rf dist/

vet:
	go vet ./...

lint:
	@which staticcheck > /dev/null 2>&1 || { echo "Install staticcheck: go install honnef.co/go/tools/cmd/staticcheck@latest"; exit 1; }
	staticcheck ./...

cross:
	mkdir -p dist
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-amd64  $(PKG)
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-arm64  $(PKG)
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-amd64   $(PKG)
	GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-arm64   $(PKG)
