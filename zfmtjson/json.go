package zfmtjson

import "encoding/json"

// Formatter encodes/decodes go struct to json format
type Formatter struct{}

// NewFormatter returns a json Formatter
func NewFormatter() Formatter {
	return Formatter{}
}

// Marshal encodes v as json
func (j Formatter) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal decodes json into v
func (j Formatter) Unmarshal(b []byte, v any) error {
	return json.Unmarshal(b, v)
}
