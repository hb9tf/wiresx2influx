package influx

import (
	"context"
	"fmt"
	"log"

	"github.com/hb9tf/wiresx2influx/wiresx"
	influxdb2 "github.com/influxdata/influxdb-client-go"
)

func Feed(ctx context.Context, logChan chan *wiresx.Log, api influxdb2.WriteApiBlocking, influxTags map[string]string) {
	for l := range logChan {
		log.Printf("%s: Message from %q (%s)\n", l.Timestamp, l.Callsign, l.Dev.InferDevice())
		var lat, lon float64
		if l.Loc != nil {
			lat = l.Loc.Lat
			lon = l.Loc.Lon
		}
		influxTags["callsign"] = l.Callsign
		influxTags["device_raw"] = string(l.Dev)
		influxTags["device"] = l.Dev.InferDevice()
		influxTags["source"] = string(l.Source)
		influxTags["location_lat"] = fmt.Sprintf("%f", lat)
		influxTags["location_lon"] = fmt.Sprintf("%f", lon)
		influxTags["description"] = l.Description
		p := influxdb2.NewPoint("callsign",
			influxTags,
			map[string]interface{}{
				"value":        l.Callsign,
				"description":  l.Description,
				"location_lat": lat,
				"location_lon": lon,
			},
			l.Timestamp)
		api.WritePoint(ctx, p)
	}
}
