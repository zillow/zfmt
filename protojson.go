package zfmt

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtoJSONFormatter encodes/decodes proto go struct to json format
type ProtoJSONFormatter struct{}

func (j *ProtoJSONFormatter) Marshall(v any) ([]byte, error) {

	if m, ok := v.(proto.Message); ok {
		return protojson.Marshal(m)
	}

	return nil, fmt.Errorf("%T, protojson formatter can only be used with proto messages", v)
}

func (j *ProtoJSONFormatter) Unmarshal(b []byte, v any) error {
	if m, ok := v.(proto.Message); ok {
		unmarshaller := protojson.UnmarshalOptions{DiscardUnknown: true}
		return unmarshaller.Unmarshal(b, m)
	}

	return fmt.Errorf("%T, protojson formatter can only be used with proto messages", v)
}
