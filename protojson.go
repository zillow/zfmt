package zfmt

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtoJSONFormatter encodes/decodes proto go struct to json format
type ProtoJSONFormatter struct{}

// Marshall ...
func (j *ProtoJSONFormatter) Marshall(v interface{}) ([]byte, error) {

	if m, ok := v.(proto.Message); ok {
		return protojson.Marshal(m)
	}

	return nil, fmt.Errorf("%T, protojson formatter can only be used with proto messages", v)
}

// Unmarshal ...
func (j *ProtoJSONFormatter) Unmarshal(b []byte, v interface{}) error {
	if m, ok := v.(proto.Message); ok {
		unmarshaller := protojson.UnmarshalOptions{DiscardUnknown: true}
		return unmarshaller.Unmarshal(b, m)
	}

	return fmt.Errorf("%T, protojson formatter can only be used with proto messages", v)
}
