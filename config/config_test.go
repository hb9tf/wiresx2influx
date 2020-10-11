package config_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/hb9tf/wiresx2influx/config"
)

func TestRead(t *testing.T) {
	tests := map[string]struct {
		inputfile string
		expected  *config.Config
	}{
		"empty config": {
			inputfile: "testdata/test0.toml",
			expected: &config.Config{
				WiresX: config.WiresX{},
				Influx: config.Influx{},
			},
		},
		"example config": {
			inputfile: "testdata/test1.toml",
			expected: &config.Config{
				WiresX: config.WiresX{
					Logfile:         "logfile.log",
					IngestWholeFile: false,
					Timezone:        "Europe/Zurich",
				},
				Influx: config.Influx{
					Server:       "http://127.0.0.1:9999",
					AuthToken:    "xyz",
					Organization: "someorg",
					Bucket:       "somebucket",
					Repeater:     "somename",
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := config.Read(tc.inputfile)
			if err != nil {
				t.Errorf("Read('%s') returned an error: %v", tc.inputfile, err)
			}
			if diff := deep.Equal(tc.expected, got); diff != nil {
				t.Errorf("Read('%s') returns unexpected result, diff: %v", tc.inputfile, diff)
			}
		})
	}
}
