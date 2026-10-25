package zfmtproto

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/zillow/zfmt/zfmtproto/testdata/example"
)

func TestRawFormatter_RoundTrip(t *testing.T) {
	// arrange
	fmtr := RawFormatter{}
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

func TestRawFormatter_Marshal_Errors(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr string
	}{
		{name: "non proto input", input: "not a proto message", wantErr: "zfmtproto: raw formatter: string is not a proto.Message"},
		{name: "invalid UTF-8 in a string field", input: invalidUTF8Example(), wantErr: "zfmtproto: encoding *example.ExampleDef"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := RawFormatter{}

			// act
			_, err := fmtr.Marshal(tt.input)

			// assert
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestRawFormatter_Unmarshal_Errors(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		target  any
		wantErr string
	}{
		{name: "non proto target", data: []byte{}, target: new(string), wantErr: "zfmtproto: raw formatter: *string is not a proto.Message"},
		{name: "typed nil target", data: []byte{}, target: (*example.ExampleDef)(nil), wantErr: "zfmtproto: raw formatter: cannot unmarshal into a nil *example.ExampleDef"},
		{name: "corrupt proto bytes", data: corruptProto, target: &example.ExampleDef{}, wantErr: "zfmtproto: decoding into *example.ExampleDef"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := RawFormatter{}

			// act
			err := fmtr.Unmarshal(tt.data, tt.target)

			// assert
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}
