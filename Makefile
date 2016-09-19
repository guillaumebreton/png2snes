all: build run

build:
	go build -o png2snes ./...

run: build
	./png2snes black.png

test: build
	go test -v ./...
