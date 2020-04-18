.PHONY: build

build:
    go get
	GOOS=windows GOARCH=amd64 go build -o wires2influx.exe wires2influx.go
