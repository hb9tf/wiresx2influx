.PHONY: build

build:
	GOOS=windows GOARCH=amd64 go build -o wires2influx.exe wires2influx.go
