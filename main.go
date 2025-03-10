package main

import (
	"fmt"

	"github.com/amirhnajafiz/prometheus-digger/internal"
)

func main() {
	// create worker pool
	pool := internal.NewWorkerPool()
	if pool == nil {
		fmt.Println("[ERR] failed to create worker pool!")
		return
	}

	// collect metrics
	pool.Collect()

	// wait for all workers
	pool.StopAndWait()
}
