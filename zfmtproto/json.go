package zfmtproto

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
)

// JSONFormatter encodes/decodes proto messages as JSON using protojson, discarding unknown fields on unmarshal.
// The zero value is ready to use.
type JSONFormatter struct{}

// Marshal encodes a proto message as json using protojson
func (JSONFormatter) Marshal(v any) ([]byte, error) {
	m, err := asMessage("json", v)
	if err != nil {
		return nil, err
	}
	b, err := protojson.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("zfmtproto: encoding %T as json: %w", v, err)
	}
	return b, nil
}

// Unmarshal decodes protojson into a proto message, discarding unknown fields
func (JSONFormatter) Unmarshal(b []byte, v any) error {
	m, err := unmarshalTarget("json", v)
	if err != nil {
		return err
	}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(b, m); err != nil {
		return fmt.Errorf("zfmtproto: decoding json into %T: %w", v, err)
	}
	return nil
}
