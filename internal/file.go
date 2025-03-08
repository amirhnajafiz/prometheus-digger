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
	fileName := metric + "_" + from + "_" + to + ".json"

	// write data to file
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
