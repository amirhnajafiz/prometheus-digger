package config

import (
	"encoding/json"

	"github.com/amirhnajafiz/prometheus-digger/pkg"
)

// LoadConfigs reads the config file and returns a config instance.
func LoadConfigs(path string) (*Config, error) {
	// read config file
	bytes, err := pkg.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// unmarshal config
	var cfg Config
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
