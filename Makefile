.PHONY: build-linux-arm64

build-linux-arm64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o makeDotApp .
