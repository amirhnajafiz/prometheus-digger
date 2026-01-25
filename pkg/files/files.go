package files

import "os"

// CheckDir checks if the input directory exists and creates it if it doesn't.
func CheckDir(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteToFile stores data into the file name.
func WriteFile(name string, data []byte) error {
	err := os.WriteFile(name, data, 0755)
	if err != nil {
		return err
	}

	return nil
}

// ReadFile reads the content of the file name.
func ReadFile(name string) ([]byte, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return data, nil
}
