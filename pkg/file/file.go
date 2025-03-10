package pkg

import "os"

// CheckDir checks if the input directory exists and creates it if it doesn't.
func CheckDir(dir string) error {
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

// WriteToFile stores data into the file name.
func WriteFile(name string, data []byte) error {
	// write data to file
	err := os.WriteFile(name, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ReadFile reads the content of the file name.
func ReadFile(name string) ([]byte, error) {
	// read data from file
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return data, nil
}
