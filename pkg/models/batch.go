package models

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Batch holds the key-values of names and queries.
type Batch struct {
	Records map[string]string
}

func NewBatch() *Batch {
	return &Batch{
		Records: make(map[string]string),
	}
}

// FillBatchFromFile accepts a file path and batch model, and fills it
// with records from the input file path.
func FillBatchFromFile(path string, b *Batch) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file `%s`: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// split by ` : `
		parts := strings.Split(line, " : ")
		if len(parts) < 2 {
			return fmt.Errorf("record `%s` is invalid. must be `name : query` format.", line)
		}

		b.Records[parts[0]] = parts[1]
	}

	return nil
}
