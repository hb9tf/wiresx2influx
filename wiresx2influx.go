package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "time/tzdata"

	"github.com/hb9tf/wiresx2influx/config"
	"github.com/hb9tf/wiresx2influx/influx"
	"github.com/hb9tf/wiresx2influx/wiresx"
	influxdb2 "github.com/influxdata/influxdb-client-go"
)

const configFileName = "wiresx2influx.conf"

func main() {
	ctx := context.Background()

	exePath, err := os.Executable()
	if err != nil {
		log.Printf("unable to determine executable path: %s", err)
		os.Exit(1)
	}
	confPath := filepath.Join(filepath.Dir(exePath), configFileName)
	log.Printf("reading config from %q", confPath)
	conf, err := config.Read(confPath)
	if err != nil {
		log.Printf("unable to load config: %s", err)
		os.Exit(1)
	}

	loc, err := time.LoadLocation(conf.WiresX.Timezone)
	if err != nil {
		log.Printf("unable to resolve location %q: %s", conf.WiresX.Timezone, err)
		os.Exit(1)
	}

	logChan := make(chan *wiresx.Log, 100)

	// Feed InfluxDB
	client := influxdb2.NewClient(conf.Influx.Server, conf.Influx.AuthToken)
	writeApi := client.WriteApiBlocking(conf.Influx.Organization, conf.Influx.Bucket)
	go func() {
		influx.Feed(ctx, logChan, writeApi, map[string]string{
			"repeater": conf.Influx.Repeater,
		})
	}()

	// Tail log
	defer close(logChan)
	if err := wiresx.TailLog(conf.WiresX.Logfile, conf.WiresX.IngestWholeFile, loc, logChan); err != nil {
		log.Printf("error in file ingestion: %s", err)
		os.Exit(1)
	}
}
