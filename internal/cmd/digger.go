package cmd

import (
	"fmt"
	"time"
)

// Digger is the main handler for fetching the metrics from Prometheus API.
type Digger struct {
	// public fields
	AddExtraCSVLabels bool
	ESC               int
	HTTPTimeout       int
	PromMetric        string
	PromURL           string
	OutputPath        string

	// private fields
	queryStep  time.Duration
	queryStart time.Time
	queryEnd   time.Time
}

// NewDigger returns an instance of digger.
func (d *Digger) Validate(start, end, step string) error {
	// convert steps to duration
	du, err := time.ParseDuration(step)
	if err != nil {
		return fmt.Errorf("invalid duration for step `%s`: %v", step, err)
	}
	d.queryStep = du

	// convert from and to into time.Time
	startDT, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return fmt.Errorf("invalid time for start `%s`: %v", start, err)
	}
	d.queryStart = startDT

	endDT, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return fmt.Errorf("invalid time for end `%s`: %v", end, err)
	}
	d.queryEnd = endDT

	// check the from and to range
	if startDT.After(endDT) {
		return fmt.Errorf("to datetime must be after from: %s - %s", start, end)
	}

	return nil
}

// Dig collects a metric from Prometheus API.
func (d *Digger) Dig() error {
	return nil
}
