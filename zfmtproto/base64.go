package zfmtproto

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/proto"
)

// Base64Formatter encodes/decodes proto messages as base64 encoded protobuf binary. Intended for use with SQS.
// The zero value is ready to use.
type Base64Formatter struct{}

// Marshal as proto and then base64 encode (useful for technologies like SQS which limit the character set)
func (Base64Formatter) Marshal(v any) ([]byte, error) {
	m, err := asMessage("base64", v)
	if err != nil {
		return nil, err
	}
	b, err := proto.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("zfmtproto: encoding %T: %w", v, err)
	}
	return base64.StdEncoding.AppendEncode(nil, b), nil
}

// Unmarshal with base64 decoding
func (Base64Formatter) Unmarshal(b []byte, v any) error {
	m, err := unmarshalTarget("base64", v)
	if err != nil {
		return err
	}
	raw, err := base64.StdEncoding.AppendDecode(nil, b)
	if err != nil {
		return fmt.Errorf("zfmtproto: decoding base64: %w", err)
	}
	if err := proto.Unmarshal(raw, m); err != nil {
		return fmt.Errorf("zfmtproto: decoding into %T: %w", v, err)
	}
	return nil
}
