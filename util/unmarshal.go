package util

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Unmarshal fills binary data into provided interface v and validates the schema ID
func UnmarshalSchema(fmtter interface {
	Marshall(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}, schemaID int, b []byte, v any) error {
	if len(b) < 5 {
		return errors.New("message does not contain schema")
	}
	schemaIDBin := b[1:5]
	id := int(binary.BigEndian.Uint32(schemaIDBin))
	// for default schema (id == 0), it is implied that the user does not care about ID and attempt to unmarshal at their own risk
	// This often happens when the topic is guaranteed to have one data type and the user would like
	// to bypass schema validation while still conforming to the confluent wire format.
	if schemaID != 0 && schemaID != id {
		return fmt.Errorf("schema IDs do not match, expect %d, got %d", schemaID, id)
	}
	return fmtter.Unmarshal(b[5:], v)
}
