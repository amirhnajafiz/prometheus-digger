package client

import (
	"math"
	"time"
)

// GetDataPoints returns the number of expected data points from the
// Prometheus API response.
func GetDataPoints(
	from,
	to time.Time,
	step time.Duration,
) int {
	fint := from.Unix()
	tint := to.Unix()
	ts := int64(math.Max(step.Seconds(), 1))

	return int((tint - fint) / ts)
}

// SplitTimeRange returns timestamps splits with maximum 11,000 expected
// datapoints from the Prometheus API response.
func SplitTimeRange(
	from,
	to time.Time,
	step time.Duration,
	maxPoints int,
) []time.Time {

	if !from.Before(to) || step <= 0 || maxPoints <= 0 {
		return nil
	}

	maxDuration := step * time.Duration(maxPoints)

	var boundaries []time.Time
	current := from
	boundaries = append(boundaries, current)

	for {
		next := current.Add(maxDuration)
		if !next.Before(to) {
			break
		}
		boundaries = append(boundaries, next)
		current = next
	}

	if boundaries[len(boundaries)-1] != to {
		boundaries = append(boundaries, to)
	}

	return boundaries
}
