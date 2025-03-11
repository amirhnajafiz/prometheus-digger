package main

import (
	"fmt"
	"log"

	"github.com/amirhnajafiz/prometheus-digger/internal/config"
	"github.com/amirhnajafiz/prometheus-digger/internal/worker"
)

func main() {
	// load configs
	cfg, err := config.LoadConfigs(configFile)
	if err != nil {
		log.Printf("[ERR] load configs failed: %v\n", err)
		return
	}

	// create worker pool
	pool := worker.NewWorkerPool(cfg)
	if pool == nil {
		fmt.Println("[ERR] failed to create worker pool!")
		return
	}

	// collect metrics
	pool.Collect()

	// wait for all workers
	pool.StopAndWait()
}
