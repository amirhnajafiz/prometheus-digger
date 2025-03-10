package internal

import "encoding/json"

// The JSON array is expected to be in the format.
type JSON struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// BytesToJSONs converts a byte array to a JSON array
// and returns an error if the conversion fails.
func BytesToJSONs(bytes []byte) ([]JSON, error) {
	var jsons []JSON

	err := json.Unmarshal(bytes, &jsons)
	if err != nil {
		return nil, err
	}

	return jsons, nil
}
