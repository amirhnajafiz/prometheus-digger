package parser

// SortDates sorts the from and to dates.
func SortDates(from, to string) (string, string) {
	// convert from and to to time.Time
	fromTime, _ := ConvertToTime(from)
	toTime, _ := ConvertToTime(to)

	// if from is after to, swap them
	if fromTime.After(toTime) {
		return to, from
	}

	return from, to
}
