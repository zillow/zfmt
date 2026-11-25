package zfmtproto

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/zillow/zfmt/zfmtproto/testdata/example"
)

func TestJSONFormatter_RoundTrip(t *testing.T) {
	// arrange
	fmtr := JSONFormatter{}
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

func TestJSONFormatter_Unmarshal_DiscardsUnknownFields(t *testing.T) {
	// arrange
	fmtr := JSONFormatter{}
	data := []byte(`{"allowed": "happy", "disallowed": 2, "MyName": "Stewart"}`)
	want := &example.ExampleDef{Allowed: "happy", Disallowed: 2}
	got := &example.ExampleDef{}

	// act
	err := fmtr.Unmarshal(data, got)

	// assert
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(want, got, protocmp.Transform()), "Unmarshal() mismatch (-want +got)")
}

func TestJSONFormatter_Marshal_Errors(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr string
	}{
		{name: "non proto input", input: map[string]string{"allowed": "happy"}, wantErr: "zfmtproto: json formatter: map[string]string is not a proto.Message"},
		{name: "invalid UTF-8 in a string field", input: invalidUTF8Example(), wantErr: "zfmtproto: encoding *example.ExampleDef as json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := JSONFormatter{}

			// act
			_, err := fmtr.Marshal(tt.input)

			// assert
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestJSONFormatter_Unmarshal_Errors(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		target  any
		wantErr string
	}{
		{name: "non proto target", data: `{}`, target: &map[string]string{}, wantErr: "zfmtproto: json formatter: *map[string]string is not a proto.Message"},
		{name: "typed nil target", data: `{}`, target: (*example.ExampleDef)(nil), wantErr: "zfmtproto: json formatter: cannot unmarshal into a nil *example.ExampleDef"},
		{name: "truncated json", data: `{"allowed": `, target: &example.ExampleDef{}, wantErr: "zfmtproto: decoding json into *example.ExampleDef"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := JSONFormatter{}

			// act
			err := fmtr.Unmarshal([]byte(tt.data), tt.target)

			// assert
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}
