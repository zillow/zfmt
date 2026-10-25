package zfmtproto

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JSONFormatter encodes/decodes proto go struct to json format
type JSONFormatter struct{}

// NewJSONFormatter returns a JSONFormatter
func NewJSONFormatter() JSONFormatter {
	return JSONFormatter{}
}

// Marshal encodes a proto message as json using protojson
func (j JSONFormatter) Marshal(v any) ([]byte, error) {
	if m, ok := v.(proto.Message); ok {
		return protojson.Marshal(m)
	}

	return nil, fmt.Errorf("%T, protojson formatter can only be used with proto messages", v)
}

// Unmarshal decodes protojson into a proto message, discarding unknown fields
func (j JSONFormatter) Unmarshal(b []byte, v any) error {
	if m, ok := v.(proto.Message); ok {
		unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
		return unmarshaler.Unmarshal(b, m)
	}

	return fmt.Errorf("%T, protojson formatter can only be used with proto messages", v)
}
