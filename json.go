package zfmt

import "encoding/json"

// JSONFormatter encodes/decodes go struct to json format
type JSONFormatter struct{}

func (j *JSONFormatter) Marshall(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JSONFormatter) Unmarshal(b []byte, v any) error {
	return json.Unmarshal(b, v)
}
