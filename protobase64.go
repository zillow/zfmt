package zfmt

import (
	"encoding/base64"
	"fmt"

	v1proto "github.com/golang/protobuf/proto"
	v2proto "google.golang.org/protobuf/proto"
)

// ProtobufBase64Formatter implements formatter interface for both protobuf v1 and v2 messages. Intended for use with SQS
type ProtobufBase64Formatter struct{}

// Marshall as proto and then base64 encode (useful for technologies like SQS which limit the character set)
func (p *ProtobufBase64Formatter) Marshall(v any) ([]byte, error) {
	switch m := v.(type) {
	case v1proto.Message:
		b, err := v1proto.Marshal(m)
		if err != nil {
			return nil, err
		}
		return []byte(base64.StdEncoding.EncodeToString(b)), nil
	case v2proto.Message:
		b, err := v2proto.Marshal(m)
		if err != nil {
			return nil, err
		}
		return []byte(base64.StdEncoding.EncodeToString(b)), nil
	default:
		return nil, fmt.Errorf("%T, proto base64 formatter can only be used with proto messages", v)
	}
}

// Unmarshal with base64 decoding
func (p *ProtobufBase64Formatter) Unmarshal(b []byte, v any) error {
	switch m := v.(type) {
	case v1proto.Message:
		raw, err := base64.StdEncoding.DecodeString(string(b))
		if err != nil {
			return err
		}
		if err := v1proto.Unmarshal(raw, m); err != nil {
			return err
		}
		return nil
	case v2proto.Message:
		raw, err := base64.StdEncoding.DecodeString(string(b))
		if err != nil {
			return err
		}
		if err = v2proto.Unmarshal(raw, m); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("%T, proto base64 formatter can only be used with proto messages", v)
	}
}
