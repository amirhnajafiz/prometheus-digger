package internal

import (
	"os"
)

const (
	// outputDir is the directory where the JSON files will be stored.
	outputDir = "output"
)

// checkDir checks if the input directory exists and creates it if it doesn't.
func checkDir(input string) error {
	// check if the input directory exists
	_, err := os.Stat(input)
	if os.IsNotExist(err) {
		// create the input directory
		err = os.Mkdir(input, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// storeMetricsInJsonFile stores the metrics in a JSON file.
func storeMetricsInJsonFile(
	metric string,
	from string,
	to string,
	data []byte,
) error {
	// check if the output directory exists
	if err := checkDir(outputDir + "/" + metric); err != nil {
		return err
	}

	// create a file name
	fileName := getFileName(metric, from, to)

	// write data to file
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// getFileName returns the file name for the given metric, from and to.
func getFileName(
	metric string,
	from string,
	to string,
) string {
	return outputDir + "/" + metric + "/" + from + "_" + to + ".json"
}
