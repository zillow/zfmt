package zfmt

import (
	"fmt"

	//nolint:staticcheck // Older zillow libs have generated code which uses this deprecated package. To maintain backwards compatability with them, the older proto serialization lib should be maintained
	v1proto "github.com/golang/protobuf/proto"
	v2proto "google.golang.org/protobuf/proto"
)

// ProtobufRawFormatter implements formatter interface for both protobuf v1 and v2 messages. Does not base64 encode.
type ProtobufRawFormatter struct{}

// Marshall encodes the data as a proto binary
func (p *ProtobufRawFormatter) Marshall(v any) ([]byte, error) {
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

// Unmarshal accepts proto binary and hydrates a proto generated struct
func (p *ProtobufRawFormatter) Unmarshal(b []byte, v any) error {
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
