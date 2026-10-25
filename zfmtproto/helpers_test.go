package zfmtproto

import (
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/zillow/zfmt/zfmtproto/testdata/example"
)

// v2OnlyMessage satisfies google.golang.org/protobuf/proto.Message without the legacy
// Reset/String/ProtoMessage methods, so it exercises the v2-only branch of each formatter.
// protoc-gen-go output implements both APIs and therefore always takes the v1 branch.
type v2OnlyMessage struct {
	msg *example.ExampleDef
}

func (m v2OnlyMessage) ProtoReflect() protoreflect.Message { return m.msg.ProtoReflect() }

func newExample() *example.ExampleDef {
	return &example.ExampleDef{
		Allowed:    "happy",
		Disallowed: 2,
		Example:    &example.ExampleDef_Foo{Foo: &example.Foo{Name: "sad"}},
	}
}

// corruptProto is a field-1 varint tag with no value, which every proto decoder rejects.
var corruptProto = []byte{0x08}

// invalidUTF8Example cannot be marshaled: proto3 string fields must hold valid UTF-8.
func invalidUTF8Example() *example.ExampleDef {
	return &example.ExampleDef{Allowed: "\xff\xfe"}
}
