.PHONY: build run clean test lint

BINARY := release-foundry
CMD    := ./cmd/release-foundry

build:
	go build -o $(BINARY) $(CMD)

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY) weekly_engineering_summary.json

test:
	go test ./...

lint:
	go vet ./...
