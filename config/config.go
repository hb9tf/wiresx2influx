// Package config provides configuration reading functionality.
package config

import (
	"fmt"
	"io"
	"os"

	"github.com/influxdata/toml"
)

// Config defines the main config structure.
type Config struct {
	WiresX WiresX `toml:"wiresx"`
	Influx Influx `toml:"influx"`
}

// WiresX defines the logfile config.
type WiresX struct {
	Logfile         string `toml:"logfile"`
	IngestWholeFile bool   `toml:"ingest_whole_file"`
	Timezone        string `toml:"timezone"`
}

// Influx defines the influx related config.
type Influx struct {
	Server       string `toml:"server"`
	AuthToken    string `toml:"auth_token"`
	Organization string `toml:"organization"`
	Bucket       string `toml:"bucket"`

	// Tags
	Repeater string `toml:"repeater"`
}

// Read parses the config file and returns the config structure.
func Read(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %q: %s", path, err)
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %q: %s", path, err)
	}
	config := &Config{}
	if err := toml.Unmarshal(buf, config); err != nil {
		return nil, fmt.Errorf("unable to unmarshal %q as TOML: %s", path, err)
	}
	return config, nil
}
