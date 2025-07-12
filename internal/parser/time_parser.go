package parser

import (
	"fmt"
	"time"
)

// ConvertSliceToTime gets a base time and an input parameters with
// this signature: +/-dd:hh:mm:ss
// Then it returns the base time adjusted by the input parameters.
func ConvertSliceToTime(base time.Time, slice string) time.Time {
	// parse the input slice
	var sign int
	if slice[0] == '-' {
		sign = -1
		slice = slice[1:]
	}
	if slice[0] == '+' {
		sign = 1
		slice = slice[1:]
	}

	// ensure the slice is in the correct format
	var days, hours, minutes, seconds int
	_, err := fmt.Sscanf(slice, "%d:%d:%d:%d", &days, &hours, &minutes, &seconds)
	if err != nil {
		return base // return the base time if parsing fails
	}

	// calculate the total duration
	duration := time.Duration(sign) * (time.Hour*24*time.Duration(days) +
		time.Hour*time.Duration(hours) +
		time.Minute*time.Duration(minutes) +
		time.Second*time.Duration(seconds))

	return base.Add(duration)
}
