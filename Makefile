
.PHONY: install dev clean build

PROG_NAME=eam  
PROG_PATH=bin/$(PROG_NAME)

install:
	go mod tidy

dev:
	DEBUG=1 go run *.go

clean:
	go clean
	rm -rf bin/
	
build: clean
	mkdir -p bin/
	go build -ldflags="-s -w" -o $(PROG_PATH) *.go
