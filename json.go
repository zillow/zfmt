package zfmt

import "encoding/json"

// JSONFormatter encodes/decodes go struct to json format
type JSONFormatter struct{}

// Marshall ...
func (j *JSONFormatter) Marshall(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal ...
func (j *JSONFormatter) Unmarshal(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}
