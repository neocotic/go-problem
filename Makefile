all: download tidy format build test bench

download:
	go mod download

tidy:
	go mod tidy -v

format:
	go fmt

build:
	go build -v ./...

test:
	go test -v ./...

bench:
	go test -run=XXX -bench=. ./...

update:
	go get -u all
