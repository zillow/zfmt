package zfmtproto

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/zillow/zfmt/zfmtproto/testdata/example"
)

func TestJSONFormatter_RoundTrip(t *testing.T) {
	// arrange
	fmtr := NewJSONFormatter()
	in := newExample()
	out := &example.ExampleDef{}

	// act
	b, err := fmtr.Marshal(in)
	require.NoError(t, err)
	err = fmtr.Unmarshal(b, out)

	// assert
	require.NoError(t, err)
	require.True(t, proto.Equal(in, out), "round trip mismatch: %v != %v", in, out)
}

func TestJSONFormatter_Unmarshal_DiscardsUnknownFields(t *testing.T) {
	// arrange
	fmtr := NewJSONFormatter()
	data := []byte(`{"allowed": "happy", "disallowed": 2, "MyName": "Stewart"}`)
	out := &example.ExampleDef{}

	// act
	err := fmtr.Unmarshal(data, out)

	// assert
	require.NoError(t, err)
	require.Equal(t, "happy", out.Allowed)
	require.Equal(t, int32(2), out.Disallowed)
}

func TestJSONFormatter_Unmarshal_InvalidJSON(t *testing.T) {
	// arrange
	fmtr := NewJSONFormatter()
	out := &example.ExampleDef{}

	// act
	err := fmtr.Unmarshal([]byte(`{"allowed": `), out)

	// assert
	require.Error(t, err)
}

func TestJSONFormatter_Marshal_RejectsNonProto(t *testing.T) {
	// arrange
	fmtr := NewJSONFormatter()

	// act
	_, err := fmtr.Marshal(map[string]string{"allowed": "happy"})

	// assert
	require.Error(t, err)
}

func TestJSONFormatter_Unmarshal_RejectsNonProto(t *testing.T) {
	// arrange
	fmtr := NewJSONFormatter()
	target := map[string]string{}

	// act
	err := fmtr.Unmarshal([]byte(`{}`), &target)

	// assert
	require.Error(t, err)
}
