default: build

all: build docker

build: build-linux

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

clean:
	rm -rf ./bin

docker:
	docker build --no-cache -t suslmk/tksinfo -f Dockerfile .