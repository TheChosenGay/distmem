build:
	go build -o bin/main.exe

run: build
	./bin/main.exe

test:
	go test -v ./...
	
.PHONY: build  test