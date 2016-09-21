all: build run

build:
	@go build -o png2snes

run: build
	@rm -rf out && mkdir out
	@./png2snes -in black.png -out-clr out/black.clr -out-pic out/black.pic

test: build
	go test -v ./...
