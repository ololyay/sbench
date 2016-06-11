VERSION=1.1.1

all:
	go build -ldflags "-s -w -X main.version=${VERSION}"
	GOOS=windows GOARCH=386 go build -ldflags "-s -w -X main.version=${VERSION}"
	GOOS=linux GOARCH=386 go build -ldflags "-s -w -X main.version=${VERSION}" -o sbench_linux
