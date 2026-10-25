package zfmtbytes

import (
	"bytes"
	"errors"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

var errBoom = errors.New("boom")

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errBoom }

func TestFormatter_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    []byte
		wantErr string
	}{
		{name: "string", input: "test", want: []byte("test")},
		{name: "byte slice", input: []byte("test"), want: []byte("test")},
		{name: "io.Reader", input: bytes.NewBufferString("test"), want: []byte("test")},
		{name: "unsupported type", input: 123, wantErr: "zfmtbytes: cannot marshal int, want string, []byte or io.Reader"},
		{name: "failing io.Reader", input: iotest.ErrReader(errBoom), wantErr: "zfmtbytes: reading"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := Formatter{}

			// act
			got, err := fmtr.Marshal(tt.input)

			// assert
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFormatter_Marshal_WrapsReaderError(t *testing.T) {
	// arrange
	fmtr := Formatter{}

	// act
	_, err := fmtr.Marshal(iotest.ErrReader(errBoom))

	// assert
	require.ErrorIs(t, err, errBoom)
}

func TestFormatter_Unmarshal_IntoString(t *testing.T) {
	// arrange
	fmtr := Formatter{}
	var got string

	// act
	err := fmtr.Unmarshal([]byte("test"), &got)

	// assert
	require.NoError(t, err)
	require.Equal(t, "test", got)
}

func TestFormatter_Unmarshal_IntoByteSliceCopiesInput(t *testing.T) {
	// arrange
	fmtr := Formatter{}
	input := []byte("test")
	got := []byte("previous contents")

	// act
	err := fmtr.Unmarshal(input, &got)
	input[0] = 'X'

	// assert
	require.NoError(t, err)
	require.Equal(t, []byte("test"), got, "target must not alias the input buffer")
}

func TestFormatter_Unmarshal_IntoWriter(t *testing.T) {
	// arrange
	fmtr := Formatter{}
	var got bytes.Buffer

	// act
	err := fmtr.Unmarshal([]byte("test"), &got)

	// assert
	require.NoError(t, err)
	require.Equal(t, "test", got.String())
}

func TestFormatter_Unmarshal_Errors(t *testing.T) {
	tests := []struct {
		name    string
		target  any
		wantErr string
	}{
		{name: "unsupported target", target: new(int), wantErr: "zfmtbytes: cannot unmarshal into *int, want *string, *[]byte or io.Writer"},
		{name: "nil *string", target: (*string)(nil), wantErr: "zfmtbytes: cannot unmarshal into a nil *string"},
		{name: "nil *[]byte", target: (*[]byte)(nil), wantErr: "zfmtbytes: cannot unmarshal into a nil *[]uint8"},
		{name: "failing io.Writer", target: failingWriter{}, wantErr: "zfmtbytes: writing to zfmtbytes.failingWriter"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := Formatter{}

			// act
			err := fmtr.Unmarshal([]byte("test"), tt.target)

			// assert
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestFormatter_Unmarshal_WrapsWriterError(t *testing.T) {
	// arrange
	fmtr := Formatter{}

	// act
	err := fmtr.Unmarshal([]byte("test"), failingWriter{})

	// assert
	require.ErrorIs(t, err, errBoom)
}
