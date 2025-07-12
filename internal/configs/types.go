package configs

import (
	"github.com/amirhnajafiz/prometheus-digger/internal/models"
)

// Config is a module that holds the configuration of the application.
type Config struct {
	URL      string         `json:"url"`
	From     string         `json:"from"`
	To       string         `json:"to"`
	PoolSize int            `json:"pool_size"`
	Queries  []models.Query `json:"queries"`
}
