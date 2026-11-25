package zfmtproto

import (
	"encoding/base64"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/zillow/zfmt/zfmtproto/testdata/example"
)

func TestBase64Formatter_RoundTrip(t *testing.T) {
	// arrange
	fmtr := Base64Formatter{}
	in := newExample()
	got := &example.ExampleDef{}

	// act
	b, err := fmtr.Marshal(in)
	require.NoError(t, err)
	err = fmtr.Unmarshal(b, got)

	// assert
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(in, got, protocmp.Transform()), "Unmarshal() mismatch (-want +got)")
}

func TestBase64Formatter_Marshal_OutputIsStdBase64(t *testing.T) {
	// arrange
	fmtr := Base64Formatter{}
	in := newExample()
	raw, err := proto.Marshal(in)
	require.NoError(t, err)

	// act
	got, err := fmtr.Marshal(in)

	// assert
	require.NoError(t, err)
	require.Equal(t, base64.StdEncoding.EncodeToString(raw), string(got))
}

func TestBase64Formatter_Marshal_Errors(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr string
	}{
		{name: "non proto input", input: 42, wantErr: "zfmtproto: base64 formatter: int is not a proto.Message"},
		{name: "invalid UTF-8 in a string field", input: invalidUTF8Example(), wantErr: "zfmtproto: encoding *example.ExampleDef"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := Base64Formatter{}

			// act
			_, err := fmtr.Marshal(tt.input)

			// assert
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestBase64Formatter_Unmarshal_Errors(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		target  any
		wantErr string
	}{
		{name: "non proto target", data: []byte("aGVsbG8="), target: new(string), wantErr: "zfmtproto: base64 formatter: *string is not a proto.Message"},
		{name: "typed nil target", data: []byte("aGVsbG8="), target: (*example.ExampleDef)(nil), wantErr: "zfmtproto: base64 formatter: cannot unmarshal into a nil *example.ExampleDef"},
		{name: "not base64", data: []byte("!!! not base64 !!!"), target: &example.ExampleDef{}, wantErr: "zfmtproto: decoding base64"},
		{name: "base64 of corrupt proto", data: []byte(base64.StdEncoding.EncodeToString(corruptProto)), target: &example.ExampleDef{}, wantErr: "zfmtproto: decoding into *example.ExampleDef"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := Base64Formatter{}

			// act
			err := fmtr.Unmarshal(tt.data, tt.target)

			// assert
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestBase64Formatter_Unmarshal_WrapsBase64Error(t *testing.T) {
	// arrange
	fmtr := Base64Formatter{}
	var corrupt base64.CorruptInputError

	// act
	err := fmtr.Unmarshal([]byte("!!!"), &example.ExampleDef{})

	// assert
	require.ErrorAs(t, err, &corrupt)
}
