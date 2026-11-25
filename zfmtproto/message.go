package zfmtproto

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// asMessage returns v as a proto.Message, or an error naming the formatter if it is not one.
func asMessage(formatter string, v any) (proto.Message, error) {
	m, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("zfmtproto: %s formatter: %T is not a proto.Message", formatter, v)
	}
	return m, nil
}

// unmarshalTarget returns v as a proto.Message that can be decoded into. A typed nil pointer such as
// (*pb.Msg)(nil) is rejected here, because protobuf dereferences it and panics.
func unmarshalTarget(formatter string, v any) (proto.Message, error) {
	m, err := asMessage(formatter, v)
	if err != nil {
		return nil, err
	}
	if !m.ProtoReflect().IsValid() {
		return nil, fmt.Errorf("zfmtproto: %s formatter: cannot unmarshal into a nil %T", formatter, v)
	}
	return m, nil
}
