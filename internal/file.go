package internal

import (
	"os"
)

// checkDir checks if the input directory exists and creates it if it doesn't.
func checkDir(dir string) error {
	// check if the input directory exists
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		// create the input directory
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// writeToFile stores data into the file name.
func writeFile(name string, data []byte) error {
	// write data to file
	err := os.WriteFile(name, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
