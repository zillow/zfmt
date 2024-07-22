package zfmt

import (
	"testing"

	"gitlab.zgtools.net/devex/archetypes/gomods/zfmt/testdata/example"
)

func TestProtoRawFormatter_MarshallUnmarshall(t *testing.T) {
	ein := example.ExampleDef{
		Allowed:    "happy",
		Disallowed: 2,
		Example: &example.ExampleDef_Foo{
			Foo: &example.Foo{
				Name: "sad",
			},
		},
	}

	fmtr := ProtobufRawFormatter{}

	b, err := fmtr.Marshall(&ein)
	if err != nil {
		t.Fatal(err)
	}

	eout := &example.ExampleDef{}
	err = fmtr.Unmarshal(b, eout)
	if err != nil {
		t.Fatal(err)
	}
	if ein.Allowed != eout.Allowed {
		t.Error("Not allowed")
	}

	if e, ok := eout.Example.(*example.ExampleDef_Foo); ok {
		if e.Foo.Name != "sad" {
			t.Error("foo sad")
		}
	} else {
		t.Error("example not foo")
	}
}
