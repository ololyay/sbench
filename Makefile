VERSION=1.0.0

all:
	go build -ldflags "-s -w -X main.version=${VERSION}"
