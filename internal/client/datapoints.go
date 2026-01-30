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
	series int,
) int {
	// series * (to - from) / step
	return int(math.Ceil(float64(to.Sub(from)/step))) * series
}

// SplitTimeRange returns timestamps splits with maximum N expected
// datapoints from the Prometheus API response.
func SplitTimeRange(
	from,
	to time.Time,
	step time.Duration,
	maxPoints int,
	totalPoints int,
) []time.Time {
	// number of requests needed
	requests := int(math.Ceil(float64(totalPoints) / float64(maxPoints)))

	// points per request
	pointsPerReq := int(math.Ceil(float64(totalPoints) / float64(requests)))

	// duration covered by each request
	chunkDuration := time.Duration(pointsPerReq-1) * step

	boundaries := make([]time.Time, 0, requests+1)
	boundaries = append(boundaries, from)

	current := from
	for range requests {
		next := current.Add(chunkDuration)
		if next.After(to) {
			next = to
		}
		boundaries = append(boundaries, next)
		current = next
	}

	return boundaries
}
