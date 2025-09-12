package schema

import (
	"bytes"
	"encoding/binary"
)

// Marshall converts input into binary data with schema ID also encoded via the wire format
func Marshall(fmtter interface {
	Marshall(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}, schemaID int, v any) ([]byte, error) {
	data, err := fmtter.Marshall(v)
	if err != nil {
		return nil, err
	}
	var body bytes.Buffer
	// version, always 0
	body.WriteByte(0)

	// 4 byte for schema ID in BigEndian
	schemaIDByte := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDByte, uint32(schemaID))
	body.Write(schemaIDByte)

	// the content
	body.Write(data)
	return body.Bytes(), nil
}
