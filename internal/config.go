package internal

import (
	"encoding/json"

	"github.com/amirhnajafiz/prometheus-digger/pkg"
)

// loadConfigs reads the config file and returns a config instance.
func loadConfigs(path string) (*config, error) {
	// read config file
	bytes, err := pkg.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// unmarshal config
	var cfg config
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
