# wiresx2influx

Parses the Wires-X log file `WiresAccess.log` and forwards the entries to a InfluxDBv2 instance.

## Config file

The config file *must* be in the same directory with the binary
and *must* be named `wiresx2influx.conf`.

The following is a copy and paste example of a config:

```
[wiresx]
logfile = "C:\\Users\\HB9TF\\OneDrive\\Documents\\WIRESXA\\AccHistory\\WiresAccess.log"
ingest_whole_file = false
timezone = "Europe/Zurich"

[influx]
server = "http://127.0.0.1:9999"
auth_token = "xyz"
organization = "someorg"
bucket = "somebucket"
repeater = "somename"  # tag
```

For `timezone`, refer to the "TZ database name" field in [this Wikipedia](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) article.

## Build

The current `Makefile` creates a Windows executable:

```
make
```

Alternatively (e.g. for development), it can be built/run with Go directly as well of course:

```
go build .
./wiresx2influx
```

Note that you can't run it with `go run .` directly as it won't find the config file that way.
