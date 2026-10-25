package zfmtproto

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/zillow/zfmt/zfmtproto/testdata/example"
)

func TestBase64Formatter_RoundTrip_V1Message(t *testing.T) {
	// arrange
	fmtr := NewBase64Formatter()
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

func TestBase64Formatter_RoundTrip_V2OnlyMessage(t *testing.T) {
	// arrange
	fmtr := NewBase64Formatter()
	in := v2OnlyMessage{msg: newExample()}
	out := v2OnlyMessage{msg: &example.ExampleDef{}}

	// act
	b, err := fmtr.Marshal(in)
	require.NoError(t, err)
	err = fmtr.Unmarshal(b, out)

	// assert
	require.NoError(t, err)
	require.True(t, proto.Equal(in.msg, out.msg), "round trip mismatch: %v != %v", in.msg, out.msg)
}

func TestBase64Formatter_Marshal_OutputIsStdBase64(t *testing.T) {
	// arrange
	fmtr := NewBase64Formatter()
	in := newExample()
	raw, err := proto.Marshal(in)
	require.NoError(t, err)

	// act
	b, err := fmtr.Marshal(in)

	// assert
	require.NoError(t, err)
	require.Equal(t, base64.StdEncoding.EncodeToString(raw), string(b))
}

func TestBase64Formatter_Marshal_RejectsNonProto(t *testing.T) {
	// arrange
	fmtr := NewBase64Formatter()

	// act
	_, err := fmtr.Marshal(42)

	// assert
	require.Error(t, err)
}

func TestBase64Formatter_Unmarshal_RejectsNonProto(t *testing.T) {
	// arrange
	fmtr := NewBase64Formatter()
	var target string

	// act
	err := fmtr.Unmarshal([]byte("aGVsbG8="), &target)

	// assert
	require.Error(t, err)
}

func TestBase64Formatter_Unmarshal_InvalidInput(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		target any
	}{
		{name: "v1 message, not base64", input: []byte("!!! not base64 !!!"), target: &example.ExampleDef{}},
		{name: "v2-only message, not base64", input: []byte("!!! not base64 !!!"), target: v2OnlyMessage{msg: &example.ExampleDef{}}},
		{name: "v1 message, base64 of corrupt proto", input: []byte(base64.StdEncoding.EncodeToString(corruptProto)), target: &example.ExampleDef{}},
		{name: "v2-only message, base64 of corrupt proto", input: []byte(base64.StdEncoding.EncodeToString(corruptProto)), target: v2OnlyMessage{msg: &example.ExampleDef{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := NewBase64Formatter()

			// act
			err := fmtr.Unmarshal(tt.input, tt.target)

			// assert
			require.Error(t, err)
		})
	}
}

func TestBase64Formatter_Marshal_PropagatesEncodingError(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{name: "v1 message", input: invalidUTF8Example()},
		{name: "v2-only message", input: v2OnlyMessage{msg: invalidUTF8Example()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := NewBase64Formatter()

			// act
			_, err := fmtr.Marshal(tt.input)

			// assert
			require.Error(t, err)
		})
	}
}
