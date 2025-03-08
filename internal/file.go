package internal

import (
	"os"
)

// storeMetricsInJsonFile stores the metrics in a JSON file.
func storeMetricsInJsonFile(
	metric string,
	from string,
	to string,
	data []byte,
) error {
	// create a file name
	fileName := getFileName(metric, from, to)

	// write data to file
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getFileName(
	metric string,
	from string,
	to string,
) string {
	// create a file name
	fileName := metric + "@" + from + "_" + to + ".json"

	return fileName
}
