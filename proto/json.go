package proto

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JSONFormatter encodes/decodes proto go struct to json format
type JSONFormatter struct{}

func NewJSONFormatter() JSONFormatter {
	return JSONFormatter{}
}

func (j *JSONFormatter) Marshall(v any) ([]byte, error) {

	if m, ok := v.(proto.Message); ok {
		return protojson.Marshal(m)
	}

	return nil, fmt.Errorf("%T, protojson formatter can only be used with proto messages", v)
}

func (j *JSONFormatter) Unmarshal(b []byte, v any) error {
	if m, ok := v.(proto.Message); ok {
		unmarshaller := protojson.UnmarshalOptions{DiscardUnknown: true}
		return unmarshaller.Unmarshal(b, m)
	}

	return fmt.Errorf("%T, protojson formatter can only be used with proto messages", v)
}
