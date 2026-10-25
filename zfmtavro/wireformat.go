package zfmtavro

import (
	"encoding/binary"
	"fmt"
)

// The Confluent wire format prefixes every payload with a 5 byte header:
// a magic byte that is always 0, followed by the schema ID as a 4 byte big-endian integer.
// https://docs.confluent.io/platform/current/schema-registry/fundamentals/serdes-develop/index.html#wire-format
const (
	magicByte = 0
	headerLen = 5
)

// encodeWireFormat prefixes payload with the Confluent wire format header for schemaID.
func encodeWireFormat(schemaID int, payload []byte) []byte {
	out := make([]byte, headerLen, headerLen+len(payload))
	out[0] = magicByte
	binary.BigEndian.PutUint32(out[1:headerLen], uint32(schemaID))
	return append(out, payload...)
}

// checkWireFormat validates the Confluent wire format header in b.
// A zero schemaID skips the schema ID check. This suits topics guaranteed to carry a single data type,
// where the caller wants the wire format without schema validation.
func checkWireFormat(schemaID int, b []byte) error {
	if len(b) < headerLen {
		return fmt.Errorf("zfmtavro: message is %d bytes, shorter than the %d byte wire format header", len(b), headerLen)
	}
	if b[0] != magicByte {
		return fmt.Errorf("zfmtavro: unknown magic byte %d, expected %d", b[0], magicByte)
	}
	id := int(binary.BigEndian.Uint32(b[1:headerLen]))
	if schemaID != 0 && schemaID != id {
		return fmt.Errorf("zfmtavro: schema ID mismatch, expected %d, got %d", schemaID, id)
	}
	return nil
}
