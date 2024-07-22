package zfmt

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/proto"
)

// ProtobufBase64Formatter implements formatter interface for both protobuf v1 and v2 messages. Intended for use with SQS
type ProtobufBase64Formatter struct{}

// Marshall with base64 encoding
func (p *ProtobufBase64Formatter) Marshall(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case proto.Message:
		b, err := proto.Marshal(m)
		if err != nil {
			return nil, err
		}
		return []byte(base64.StdEncoding.EncodeToString(b)), nil
	default:
		return nil, fmt.Errorf("%T, proto base64 formatter can only be used with proto messages", v)
	}
}

// Unmarshal with base64 decoding
func (p *ProtobufBase64Formatter) Unmarshal(b []byte, v interface{}) error {
	switch m := v.(type) {
	case proto.Message:
		raw, err := base64.StdEncoding.DecodeString(string(b))
		if err != nil {
			return err
		}
		if err = proto.Unmarshal(raw, m); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("%T, proto base64 formatter can only be used with proto messages", v)
	}
}
