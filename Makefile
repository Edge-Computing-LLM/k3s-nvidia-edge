BINARY := k3s-nvidia-edge

.PHONY: build test fmt vet clean
.PHONY: check install-local

build:
	go build -o bin/$(BINARY) ./cmd/k3s-nvidia-edge

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

vet:
	go vet ./...

check: fmt vet test build

install-local: build
	install -D -m 0755 bin/$(BINARY) $(HOME)/.local/bin/$(BINARY)

clean:
	rm -rf bin
