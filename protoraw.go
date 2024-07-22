package zfmt

import (
	"fmt"

	v1proto "github.com/golang/protobuf/proto"
	v2proto "google.golang.org/protobuf/proto"
)

// ProtobufRawFormatter implements formatter interface for both protobuf v1 and v2 messages. Does not base64 encode.
type ProtobufRawFormatter struct{}

// Marshall ...
// same as proto.go formatter but does not base64 encode messages
func (p *ProtobufRawFormatter) Marshall(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case v1proto.Message:
		b, err := v1proto.Marshal(m)
		if err != nil {
			return nil, err
		}
		return b, nil
	case v2proto.Message:
		b, err := v2proto.Marshal(m)
		if err != nil {
			return nil, err
		}
		return b, nil
	default:
		return nil, fmt.Errorf("%T, protoraw formatter can only be used with proto messages", v)
	}
}

// Unmarshal ...
func (p *ProtobufRawFormatter) Unmarshal(b []byte, v interface{}) error {
	switch m := v.(type) {
	case v1proto.Message:
		if err := v1proto.Unmarshal(b, m); err != nil {
			return err
		}
		return nil
	case v2proto.Message:
		if err := v2proto.Unmarshal(b, m); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("%T, protoraw formatter can only be used with proto messages", v)
	}
}
