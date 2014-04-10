all:
	go get -u ./...
	go build -a -v -race

clean:
	go clean
