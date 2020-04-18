package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/influxdata/toml"
)

type Config struct {
	Repeater string `toml:"repeater"`

	WiresX WiresX `toml:"wiresx"`
	Influx Influx `toml:"influx"`
	Slack  Slack  `toml:"slack"`
}

type WiresX struct {
	Logfile         string `toml:"logfile"`
	IngestWholeFile bool   `toml:"ingest_whole_file"`
	Timezone        string `toml:"timezone"`
}

type Influx struct {
	Server       string `toml:"server"`
	AuthToken    string `toml:"auth_token"`
	Organization string `toml:"organization"`
	Bucket       string `toml:"bucket"`
	Dry          bool   `toml:"dry"`
}

type Slack struct {
	Webhook string `toml:"webhook"`
	Dry     bool   `toml:"dry"`
}

func Read(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %q: %s", path, err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %q: %s", path, err)
	}
	config := &Config{}
	if err := toml.Unmarshal(buf, config); err != nil {
		return nil, fmt.Errorf("unable to unmarshal %q as TOML: %s", path, err)
	}
	return config, nil
}
