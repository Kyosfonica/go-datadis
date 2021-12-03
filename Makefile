bin:
	mkdir -p bin
	go build -o bin/datadis ./cmd/datadis

.PHONY: bin
