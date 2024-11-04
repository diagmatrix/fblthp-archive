all: run

build:
	go build -o bin/fblthp-archive

run: build
	./bin/fblthp-archive

test:
	go test -v ./...