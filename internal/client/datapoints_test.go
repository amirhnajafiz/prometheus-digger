package client

import (
	"testing"
	"time"
)

func TestGetDataPoints(t *testing.T) {
	cases := [][]time.Time{
		{
			time.Now(), time.Now().Add(10 * time.Second),
			time.Now(), time.Now().Add(10 * time.Second),
			time.Now(), time.Now().Add(1 * time.Hour),
			time.Now(), time.Now().Add(1 * time.Hour),
		},
	}

	steps := []time.Duration{
		time.Second,
		5 * time.Second,
		time.Second,
		5 * time.Second,
	}

	want := []int{
		10,
		2,
		3600,
		720,
	}

	for i := range len(cases) {
		got := GetDataPoints(cases[i][0], cases[i][1], steps[i], 1)
		if got != want[i] {
			t.Errorf(
				"input: (%s, %s, %f), output: %d, expected: %d",
				cases[i][0].Format(time.RFC3339),
				cases[i][1].Format(time.RFC3339),
				steps[i].Seconds(),
				got,
				want[i],
			)
		}
	}
}

func TestSplitTimeRange(t *testing.T) {
	cases := [][]time.Time{
		{
			time.Now(), time.Now().Add(10 * time.Second),
			time.Now(), time.Now().Add(10 * time.Second),
			time.Now(), time.Now().Add(1 * time.Hour),
			time.Now(), time.Now().Add(1 * time.Hour),
		},
	}

	steps := []time.Duration{
		time.Second,
		5 * time.Second,
		time.Second,
		5 * time.Second,
	}

	maxPoints := []int{
		2,
		2,
		20,
		20,
	}

	want := []int{
		7,
		4,
		182,
		126,
	}

	for i := range len(cases) {
		got := SplitTimeRange(cases[i][0], cases[i][1], steps[i], maxPoints[i])
		if len(got) != want[i] {
			t.Errorf(
				"input: (%s, %s, %f, %d), output: %d, expected: %d",
				cases[i][0].Format(time.RFC3339),
				cases[i][1].Format(time.RFC3339),
				steps[i].Seconds(),
				maxPoints[i],
				len(got),
				want[i],
			)
		}
	}
}
