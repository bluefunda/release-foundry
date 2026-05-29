.PHONY: build run clean test race lint vet install snapshot

BINARY := release-foundry
CMD    := ./cmd/release-foundry
VERSION ?= dev

build:
	go build -ldflags "-X main.buildVersion=$(VERSION)" -o $(BINARY) $(CMD)

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY) release-summary.json

test:
	go test ./...

race:
	go test -race ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

install:
	go install -ldflags "-X main.buildVersion=$(VERSION)" $(CMD)

snapshot:
	goreleaser release --snapshot --clean
