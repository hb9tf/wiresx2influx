package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"4d63.com/tz"
	"github.com/hb9tf/wiresx2influx/config"
	"github.com/hb9tf/wiresx2influx/influx"
	"github.com/hb9tf/wiresx2influx/slack"
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

	loc, err := tz.LoadLocation(conf.WiresX.Timezone)
	if err != nil {
		log.Printf("unable to resolve location %q: %s", conf.WiresX.Timezone, err)
		os.Exit(1)
	}

	chans := []chan *wiresx.Log{}
	tags := map[string]string{
		"repeater": conf.Repeater,
	}

	// Feed InfluxDB
	if conf.Influx.Server != "" {
		log.Printf("starting InfluxDB feeder to %q", conf.Influx.Server)
		influxLogChan := make(chan *wiresx.Log, 100)
		chans = append(chans, influxLogChan)

		client := influxdb2.NewClient(conf.Influx.Server, conf.Influx.AuthToken)
		writeApi := client.WriteApiBlocking(conf.Influx.Organization, conf.Influx.Bucket)
		go func() {
			influx.Feed(ctx, influxLogChan, writeApi, tags, conf.Influx.Dry)
		}()
	}

	// Feed Slack
	if conf.Slack.Webhook != "" {
		log.Printf("starting Slack feeder")
		slackLogChan := make(chan *wiresx.Log, 100)
		chans = append(chans, slackLogChan)

		slackClient := &slack.Slacker{
			conf.Slack.Webhook,
			&http.Client{},
		}
		go func() {
			slack.Feed(ctx, slackLogChan, slackClient, tags, conf.Slack.Dry)
		}()
	}

	// Tail log
	if err := wiresx.TailLog(conf.WiresX.Logfile, conf.WiresX.IngestWholeFile, loc, chans); err != nil {
		log.Printf("error in file ingestion: %s", err)
		os.Exit(1)
	}
}
