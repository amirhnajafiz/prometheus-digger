package internal

import (
	"os"
)

const (
	outputDir = "output"
)

// checkOutputDir checks if the output directory exists and creates it if it doesn't.
func checkOutputDir() error {
	// check if the output directory exists
	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		// create the output directory
		err = os.Mkdir(outputDir, 0755)
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
	fileName := outputDir + "/" + metric + "@" + from + "_" + to + ".json"

	return fileName
}
