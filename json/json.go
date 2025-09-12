package json

import "encoding/json"

// Formatter encodes/decodes go struct to json format
type Formatter struct{}

func NewFormatter() Formatter {
	return Formatter{}
}

func (j *Formatter) Marshall(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j *Formatter) Unmarshal(b []byte, v any) error {
	return json.Unmarshal(b, v)
}
