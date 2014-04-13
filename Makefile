all:
	go get -u ./...
	go build -a -v -race
	go install

install:
	go install

dep:
	go get -u ./...

build:
	go build -a -v -race

clean:
	go clean
