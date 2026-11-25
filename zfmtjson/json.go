package zfmtjson

import (
	jsonv1 "encoding/json"
	"encoding/json/v2"
)

// Formatter encodes/decodes Go values as JSON using encoding/json/v2.
//
// It keeps the wire behavior of encoding/json (v1) wherever v2's defaults would change it silently:
// nil slices and maps marshal as null, and object names match struct fields case-insensitively
// (without ignoring '_' or '-', as in v1). v2's stricter checks stay on, so invalid UTF-8 and
// duplicate object names are errors rather than being repaired or overwritten.
//
// The zero value is ready to use.
type Formatter struct{}

var (
	marshalOptions = json.JoinOptions(
		json.FormatNilSliceAsNull(true),
		json.FormatNilMapAsNull(true),
	)
	unmarshalOptions = json.JoinOptions(
		json.MatchCaseInsensitiveNames(true),
		jsonv1.MatchCaseSensitiveDelimiter(true),
	)
)

// Marshal encodes v as json
func (Formatter) Marshal(v any) ([]byte, error) {
	return json.Marshal(v, marshalOptions)
}

// Unmarshal decodes json into v
func (Formatter) Unmarshal(b []byte, v any) error {
	return json.Unmarshal(b, v, unmarshalOptions)
}
