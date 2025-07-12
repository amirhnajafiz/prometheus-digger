package configs

import (
	"encoding/json"

	"github.com/amirhnajafiz/prometheus-digger/internal/models"
	"github.com/amirhnajafiz/prometheus-digger/pkg"
)

// Config is a module that holds the configuration of the application.
type Config struct {
	URL      string         `json:"url"`
	From     string         `json:"from"`
	To       string         `json:"to"`
	PoolSize int            `json:"pool_size"`
	Queries  []models.Query `json:"queries"`
}

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
