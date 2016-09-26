all: build run

build:
	@go build -o png2snes

run: build
	@rm -rf out && mkdir out
	@./png2snes -in examples/test.tmx

test: build
	go test -v ./...
