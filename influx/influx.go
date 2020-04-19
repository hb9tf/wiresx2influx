package influx

import (
	"context"
	"log"

	"github.com/hb9tf/wiresx2influx/wiresx"
	influxdb2 "github.com/influxdata/influxdb-client-go"
)

func Feed(ctx context.Context, logChan chan *wiresx.Log, api influxdb2.WriteApiBlocking, tags map[string]string, dry bool) {
	for l := range logChan {
		var lat, lon float64
		if l.Loc != nil {
			lat = l.Loc.Latitude
			lon = l.Loc.Longitude
		}
		p := influxdb2.NewPoint("callsign",
			tags,
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
		if dry {
			log.Printf("DRY: Sending message to InfluxDB: %+v", p)
			continue
		}
		api.WritePoint(ctx, p)
	}
}
