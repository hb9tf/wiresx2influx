package wiresx

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hb9tf/wiresx2influx/geo"
	"github.com/hpcloud/tail"
)

const (
	logTimeFmt = "2006/01/02 15:04:05"
)

var (
	deviceTypes = map[string]string{
		"E0": "FT-1D",
		"E5": "FT-2D",
		"EA": "FT-3D",
		"F0": "FTM-400D",
		"F5": "FTM-100D",
		"FA": "FTM-300D",
		"G0": "FT-991",
		"H0": "FTM-3200D",
		"H5": "FT-70D",
		"HA": "FTM-3207D",
		"HF": "FTM-7250D",
		"R":  "repeater",
	}
)

type Activity string

type Device string

func (d Device) InferDevice() string {
	dev := strings.ToUpper(string(d))
	for k, v := range deviceTypes {
		if strings.HasPrefix(dev, k) {
			return v
		}
	}
	id, err := strconv.ParseInt(dev, 10, 32)
	if err != nil {
		return "unknown"
	}
	if id/10000%2 == 0 {
		return "room"
	}
	return "node"
}

// https://github.com/HB9UF/unconfusion
type Log struct {
	Callsign    string
	Dev         Device
	Description string
	Timestamp   time.Time
	Source      Activity
	Loc         *geo.Location
}

// parseLocation attempts to parse the WiresX log location format.
// Examples:
// N:36 48' 58" / W:084 10' 02"
// Lat:N:46 01' 00" / Lon:E:007 44' 36" / R:001km /
func parseLocation(l string) (*geo.Location, error) {
	l = strings.TrimSpace(l)
	if l == "" {
		return nil, nil
	}
	l = strings.ReplaceAll(l, "Lat:", "")
	l = strings.ReplaceAll(l, "Lon:", "")

	var err error
	var lat, lon float64
	parts := strings.Split(l, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("unexpected location tokens: %s", l)
	}

	latParts := strings.Split(strings.TrimSpace(parts[0]), ":")
	latClean := strings.ReplaceAll(latParts[1], "'", "")
	latClean = strings.ReplaceAll(latClean, "\"", "")
	latDegParts := strings.Split(latClean, " ")
	lat, err = strconv.ParseFloat(latDegParts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("conversion to location failed (lat deg): %s", l)
	}
	latMin, err := strconv.ParseFloat(latDegParts[1], 64)
	if err != nil {
		return nil, fmt.Errorf("conversion to location failed (lat min): %s", l)
	}
	latSec, err := strconv.ParseFloat(latDegParts[2], 64)
	if err != nil {
		return nil, fmt.Errorf("conversion to location failed (lat sec): %s", l)
	}
	lat += ((latSec / 60) + latMin) / 60
	if latParts[0] == "S" {
		lat = lat * -1
	}

	lonParts := strings.Split(strings.TrimSpace(parts[1]), ":")
	lonClean := strings.ReplaceAll(lonParts[1], "'", "")
	lonClean = strings.ReplaceAll(lonClean, "\"", "")
	lonDegParts := strings.Split(lonClean, " ")
	lon, err = strconv.ParseFloat(lonDegParts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("conversion to location failed (lon deg): %s", l)
	}
	lonMin, err := strconv.ParseFloat(lonDegParts[1], 64)
	if err != nil {
		return nil, fmt.Errorf("conversion to location failed (lon min): %s", l)
	}
	lonSec, err := strconv.ParseFloat(lonDegParts[2], 64)
	if err != nil {
		return nil, fmt.Errorf("conversion to location failed (lon sec): %s", l)
	}
	lon += ((lonSec / 60) + lonMin) / 60
	if lonParts[0] == "W" {
		lon = lon * -1
	}

	loc, err := geo.Lookup(lat, lon)
	if err != nil {
		return &geo.Location{Latitude: lat, Longitude: lon}, nil
	}
	return loc, nil
}

func parseLogline(line string, timeLoc *time.Location) (*Log, error) {
	parts := strings.Split(line, "%")
	if len(parts) != 13 {
		return nil, fmt.Errorf("unexpected amount of tokens (want: 13): %s", line)
	}
	ts, err := time.ParseInLocation(logTimeFmt, parts[3], timeLoc)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp format (want: %q): %s", logTimeFmt, parts[3])
	}
	loc, err := parseLocation(parts[6])
	if err != nil {
		log.Printf("invalid location: %s", parts[6])
	}
	return &Log{
		Callsign:    strings.ToUpper(parts[0]),
		Dev:         Device(strings.ToUpper(parts[1])),
		Description: parts[2],
		Timestamp:   ts,
		Source:      Activity(parts[4]),
		Loc:         loc,
	}, nil
}

func TailLog(path string, ingestWholeFile bool, loc *time.Location, chans []chan *Log) error {
	whence := os.SEEK_END
	if ingestWholeFile {
		whence = os.SEEK_SET
	}
	t, err := tail.TailFile(path, tail.Config{
		ReOpen:    true,  // Reopen recreated files (tail -F)
		MustExist: true,  // Fail early if the file does not exist
		Poll:      false, // Poll for file changes instead of using inotify
		Follow:    true,  // Continue looking for new lines (tail -f)
		// Logger, when nil, is set to tail.DefaultLogger
		// To disable logging: set field to tail.DiscardingLogger
		Logger: tail.DiscardingLogger,
		Location: &tail.SeekInfo{
			Whence: whence,
			Offset: 0,
		},
	})
	if err != nil {
		return fmt.Errorf("unable to tail file %q: %s", path, err)
	}
	for line := range t.Lines {
		wl, err := parseLogline(line.Text, loc)
		if err != nil {
			log.Printf("error parsing log line: %s", err)
			continue
		}
		log.Printf("%s: Message from %q (%s)\n", wl.Timestamp, wl.Callsign, wl.Dev.InferDevice())
		for _, c := range chans {
			c <- wl
		}
	}
	return nil
}
