package influx

import (
	"context"
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
		p := influxdb2.NewPoint("callsign",
			influxTags,
			map[string]interface{}{
				"value":        l.Callsign,
				"device_raw":   string(l.Dev),
				"device":       l.Dev.InferDevice(),
				"source":       string(l.Source),
				"location_lat": lat,
				"location_lon": lon,
				"description":  l.Description,
			},
			l.Timestamp)
		api.WritePoint(ctx, p)
	}
}
