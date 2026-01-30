package configs

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/amirhnajafiz/prometheus-digger/pkg/files"
)

// Config holds the configuration of the the digger.
type Config struct {
	AddExtraCSVLabels    bool   `json:"extra_csv_labels"`
	EstimatedSeriesCount int    `json:"esc"`
	RequestTimeout       int    `json:"request_timeout"`
	DataDir              string `json:"data_directory"`
	PrometheusURL        string `json:"prometheus_url"`
	Step                 string `json:"step"`
}

// LoadConfigs reads a json format config file and returns a config instance.
func LoadConfigs(path string) (*Config, error) {
	// read config file
	bytes, err := files.ReadFile(path)
	if err != nil {
		log.Printf("failed to read config file `%s`: %v", path, err)
	}

	// unmarshal configs into an instance
	cfg := Default()

	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config file `%s`: %v", path, err)
		}
	}

	// trim '/' from Prometheus URL
	cfg.PrometheusURL = strings.Trim(cfg.PrometheusURL, "/")

	return &cfg, nil
}
