package zfmt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// SchematizedAvroFormatter follows the Confluent Wire Format https://docs.confluent.io/current/schema-registry/serdes-develop/index.html#wire-format
type SchematizedAvroFormatter struct {
	formatter AvroFormatter
	SchemaID  int
}

// Marshall converts input into avro binary data with schema ID attached
func (p *SchematizedAvroFormatter) Marshall(v interface{}) ([]byte, error) {
	return marshall(&p.formatter, p.SchemaID, v)
}

// Unmarshal fills avro binary data into provided interface v and validates the schema ID
func (p *SchematizedAvroFormatter) Unmarshal(b []byte, v interface{}) error {
	return unmarshal(&p.formatter, p.SchemaID, b, v)
}

// marshall converts input into binary data with schema ID also encoded via the wire format
func marshall(fmtter Formatter, schemaID int, v interface{}) ([]byte, error) {
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

// Unmarshal fills binary data into provided interface v and validates the schema ID
func unmarshal(fmtter Formatter, schemaID int, b []byte, v interface{}) error {
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
