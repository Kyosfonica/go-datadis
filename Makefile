bin:
	mkdir -p bin
	env GOOS=windows GOARCH=amd64 go build -o bin/datadis.exe ./cmd/datadis

.PHONY: bin
