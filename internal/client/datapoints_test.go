package client

import (
	"testing"
	"time"
)

func TestGetDataPoints(t *testing.T) {
	// from and to
	cases := [][]time.Time{
		{time.Now(), time.Now().Add(10 * time.Second)},
		{time.Now(), time.Now().Add(10 * time.Second)},
		{time.Now(), time.Now().Add(10 * time.Minute)},
		{time.Now(), time.Now().Add(30 * time.Minute)},
		{time.Now(), time.Now().Add(1 * time.Hour)},
		{time.Now(), time.Now().Add(1 * time.Hour)},
	}

	// for each case there is a step
	steps := []time.Duration{
		time.Second,
		5 * time.Second,
		time.Second,
		5 * time.Second,
		time.Second,
		5 * time.Second,
	}

	// for each case there is number of data points
	want := []int{
		10,
		2,
		600,
		360,
		3600,
		720,
	}

	// loop over and submit cases
	for i := range len(cases) {
		got := GetDataPoints(cases[i][0], cases[i][1], steps[i], 1)
		if got != want[i] {
			t.Errorf(
				"[%d] input: (%s, %s, %f), output: %d, expected: %d",
				i,
				cases[i][0].Format(time.RFC3339),
				cases[i][1].Format(time.RFC3339),
				steps[i].Seconds(),
				got,
				want[i],
			)
		}
	}

	t.Logf("passed %d cases", len(cases))
}

func TestSplitTimeRange(t *testing.T) {
	// from and to
	cases := [][]time.Time{
		{time.Now(), time.Now().Add(10 * time.Second)},
		{time.Now(), time.Now().Add(10 * time.Second)},
		{time.Now(), time.Now().Add(10 * time.Minute)},
		{time.Now(), time.Now().Add(30 * time.Minute)},
		{time.Now(), time.Now().Add(1 * time.Hour)},
		{time.Now(), time.Now().Add(1 * time.Hour)},
	}

	// for each case there is a step
	steps := []time.Duration{
		time.Second,
		5 * time.Second,
		time.Second,
		5 * time.Second,
		time.Second,
		5 * time.Second,
	}

	// for each case there is max points within each range
	maxPoints := []int{
		2,
		2,
		2,
		20,
		20,
		20,
	}

	// for each case there is total data points
	dp := []int{
		10,
		2,
		600,
		360,
		3600,
		720,
	}

	// for each case we need this number of timestamps
	// time ranges size = (totalDatapoints/maxPoints) + 1
	// time ranges step = (to - from) / (time ranges size - 1)
	want := []int{
		6,
		2,
		301,
		19,
		181,
		37,
	}

	for i := range len(cases) {
		got := SplitTimeRange(cases[i][0], cases[i][1], steps[i], maxPoints[i], dp[i])
		if len(got) != want[i] {
			t.Errorf(
				"[%d] input: (%s, %s, %f, %d, %d), output: %d, expected: %d",
				i,
				cases[i][0].Format(time.RFC3339),
				cases[i][1].Format(time.RFC3339),
				steps[i].Seconds(),
				maxPoints[i],
				dp[i],
				len(got),
				want[i],
			)
		}
	}

	t.Logf("passed %d cases", len(cases))
}
