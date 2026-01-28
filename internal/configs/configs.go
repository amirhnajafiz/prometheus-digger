package configs

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/amirhnajafiz/prometheus-digger/pkg/files"
)

// Config holds the configuration of the the digger.
type Config struct {
	RequestTimeout int    `json:"request_timeout"`
	DataDir        string `json:"data_directory"`
	PrometheusURL  string `json:"prometheus_url"`
	Steps          string `json:"steps"`
}

// LoadConfigs reads a json format config file and returns a config instance.
func LoadConfigs(path string) (*Config, error) {
	// read config file
	bytes, err := files.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file `%s`: %v", path, err)
	}

	// unmarshal configs into an instance
	cfg := Default()
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file `%s`: %v", path, err)
	}

	// trim '/' from Prometheus URL
	cfg.PrometheusURL = strings.Trim(cfg.PrometheusURL, "/")

	return &cfg, nil
}
