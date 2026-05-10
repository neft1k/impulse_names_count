.PHONY: build test bench vet lint run clean

BINARY := count-names

build:
	go build -o $(BINARY) ./cmd/count-names

test:
	go test ./...

bench:
	go test -bench=. -benchmem ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

run: build
	./$(BINARY) $(ARGS)

clean:
	rm -f $(BINARY)
