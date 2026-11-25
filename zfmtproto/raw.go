package zfmtproto

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// RawFormatter encodes/decodes proto messages using the protobuf binary wire format. It does not base64 encode.
// The zero value is ready to use.
type RawFormatter struct{}

// Marshal encodes the data as a proto binary
func (RawFormatter) Marshal(v any) ([]byte, error) {
	m, err := asMessage("raw", v)
	if err != nil {
		return nil, err
	}
	b, err := proto.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("zfmtproto: encoding %T: %w", v, err)
	}
	return b, nil
}

// Unmarshal accepts proto binary and hydrates a proto generated struct
func (RawFormatter) Unmarshal(b []byte, v any) error {
	m, err := unmarshalTarget("raw", v)
	if err != nil {
		return err
	}
	if err := proto.Unmarshal(b, m); err != nil {
		return fmt.Errorf("zfmtproto: decoding into %T: %w", v, err)
	}
	return nil
}
