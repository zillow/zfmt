package zfmtproto

import (
	"github.com/zillow/zfmt/zfmtproto/testdata/example"
)

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
