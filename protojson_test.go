package zfmt

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.zgtools.net/devex/archetypes/gomods/zfmt/testdata/example"
)

func TestProtoJSONFormatter_MarshallUnmarshall(t *testing.T) {
	ein := example.ExampleDef{
		Allowed:    "happy",
		Disallowed: 2,
		Example: &example.ExampleDef_Foo{
			Foo: &example.Foo{
				Name: "sad",
			},
		},
	}

	fmtr := ProtoJSONFormatter{}

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

func TestProtoJSONFormatter_UnmarshallWithUnknown(t *testing.T) {
	data := "{\n    \"allowed\": \"happy\",\n    \"disallowed\": 2,\n    \"MyName\": \"Stewart\"\n}"

	fmtr := ProtoJSONFormatter{}
	eout := &example.ExampleDef{}
	err := fmtr.Unmarshal([]byte(data), eout)
	require.NoError(t, err)
	require.Equal(t, "happy", eout.Allowed)
	require.Equal(t, int32(2), eout.Disallowed)
}
