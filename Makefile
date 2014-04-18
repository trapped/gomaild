all:
	go get -u ./...
	go build -a -v -race

all-norace:
	go get -u ./...
	go build -a -v

install:
	go install

deps:
	go get -u ./...

build:
	go build -a -v -race

build-norace:
	go build -a -v

clean:
	go clean
