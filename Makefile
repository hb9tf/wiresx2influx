.PHONY: build

build:
	go get
	GOOS=windows GOARCH=amd64 go build -o wiresx2influx.exe wiresx2influx.go

clean:
	rm wiresx2influx.exe
