package zfmtproto

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/zillow/zfmt/zfmtproto/testdata/example"
)

func TestRawFormatter_RoundTrip_V1Message(t *testing.T) {
	// arrange
	fmtr := NewRawFormatter()
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

func TestRawFormatter_RoundTrip_V2OnlyMessage(t *testing.T) {
	// arrange
	fmtr := NewRawFormatter()
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

func TestRawFormatter_Marshal_RejectsNonProto(t *testing.T) {
	// arrange
	fmtr := NewRawFormatter()

	// act
	_, err := fmtr.Marshal("not a proto message")

	// assert
	require.Error(t, err)
}

func TestRawFormatter_Unmarshal_RejectsNonProto(t *testing.T) {
	// arrange
	fmtr := NewRawFormatter()
	var target string

	// act
	err := fmtr.Unmarshal([]byte{}, &target)

	// assert
	require.Error(t, err)
}

func TestRawFormatter_Unmarshal_InvalidBytes(t *testing.T) {
	tests := []struct {
		name   string
		target any
	}{
		{name: "v1 message", target: &example.ExampleDef{}},
		{name: "v2-only message", target: v2OnlyMessage{msg: &example.ExampleDef{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := NewRawFormatter()

			// act
			err := fmtr.Unmarshal(corruptProto, tt.target)

			// assert
			require.Error(t, err)
		})
	}
}

func TestRawFormatter_Marshal_PropagatesEncodingError(t *testing.T) {
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
			fmtr := NewRawFormatter()

			// act
			_, err := fmtr.Marshal(tt.input)

			// assert
			require.Error(t, err)
		})
	}
}
