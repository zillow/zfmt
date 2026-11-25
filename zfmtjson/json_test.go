package zfmtjson

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"github.com/zillow/zfmt/zfmtjson/testdata"
)

func TestFormatter_RoundTrip(t *testing.T) {
	// arrange
	fmtr := Formatter{}
	in := testdata.ExampleJson{Id: "hello"}
	got := testdata.ExampleJson{}

	// act
	b, err := fmtr.Marshal(in)
	require.NoError(t, err)
	err = fmtr.Unmarshal(b, &got)

	// assert
	require.NoError(t, err)
	require.JSONEq(t, `{"Id":"hello"}`, string(b))
	require.Empty(t, cmp.Diff(in, got), "Unmarshal() mismatch (-want +got)")
}

// The tests below pin the wire behavior the formatter keeps from encoding/json (v1), and the
// stricter v2 checks it keeps on, so a change of JSON package or options shows up as a test failure.

func TestFormatter_Marshal_NilSliceAndMapAreNull(t *testing.T) {
	// arrange
	fmtr := Formatter{}
	in := struct {
		Tags   []string
		Labels map[string]string
	}{}

	// act
	b, err := fmtr.Marshal(in)

	// assert
	require.NoError(t, err)
	require.JSONEq(t, `{"Tags":null,"Labels":null}`, string(b))
}

func TestFormatter_Unmarshal_NameMatching(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		wantID string
	}{
		{name: "exact name", data: `{"Id":"x"}`, wantID: "x"},
		{name: "different case matches, as in v1", data: `{"id":"x"}`, wantID: "x"},
		{name: "underscore is not ignored, as in v1", data: `{"i_d":"x"}`, wantID: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := Formatter{}
			got := testdata.ExampleJson{}

			// act
			err := fmtr.Unmarshal([]byte(tt.data), &got)

			// assert
			require.NoError(t, err)
			require.Equal(t, tt.wantID, got.Id)
		})
	}
}

// json/v2 randomizes parts of its error text (for example "cannot" vs "unable to"), so wantErr
// matches only the stable tail of each message.
func TestFormatter_Marshal_Errors(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr string
	}{
		{name: "unsupported type", input: make(chan int), wantErr: "marshal from Go chan int"},
		{name: "invalid UTF-8", input: "\xff", wantErr: "invalid UTF-8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := Formatter{}

			// act
			_, err := fmtr.Marshal(tt.input)

			// assert
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestFormatter_Unmarshal_Errors(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr string
	}{
		{name: "truncated json", data: `{"Id":`, wantErr: "unexpected EOF"},
		{name: "duplicate object name", data: `{"Id":"a","Id":"b"}`, wantErr: "duplicate object member name"},
		{name: "duplicate name differing only in case", data: `{"Id":"a","id":"b"}`, wantErr: "duplicate object member name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := Formatter{}
			got := testdata.ExampleJson{}

			// act
			err := fmtr.Unmarshal([]byte(tt.data), &got)

			// assert
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}
