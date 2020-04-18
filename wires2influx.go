package main

import (
	"context"
	"log"
	"os"

	"4d63.com/tz"
	"github.com/hb9tf/wiresx2influx/influx"
	"github.com/hb9tf/wiresx2influx/wiresx"
	influxdb2 "github.com/influxdata/influxdb-client-go"
)

func main() {
	ctx := context.Background()

	// TODO: Needs to move into a config file.
	infile := "WiresAccess.log"
	ingestWholeFile := false
	timezone := "Europe/Zurich"
	influxServer := "http://192.168.73.12:9999"
	influxAuth := ""
	influxOrg := "hb9tf"
	influxBucket := "wiresx"
	influxTags := map[string]string{
		"relais":        "lszh",
		"wiresx2influx": "0.1",
	}

	loc, err := tz.LoadLocation(timezone)
	if err != nil {
		log.Printf("unable to resolve location %q: %s", timezone, err)
		os.Exit(1)
	}

	logChan := make(chan *wiresx.Log, 100)

	// Feed InfluxDB
	client := influxdb2.NewClient(influxServer, influxAuth)
	writeApi := client.WriteApiBlocking(influxOrg, influxBucket)
	go func() {
		influx.Feed(ctx, logChan, writeApi, influxTags)
	}()

	// Tail log
	defer close(logChan)
	if err := wiresx.TailLog(infile, ingestWholeFile, loc, logChan); err != nil {
		log.Printf("error in file ingestion: %s", err)
		os.Exit(1)
	}
}
