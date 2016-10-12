all: build run

build:
	@go build -o tmx2snes

run: build
	@rm -rf out && mkdir out
	@./tmx2snes -in resources/map.tmx -out out

test: build
	go test -v ./...
