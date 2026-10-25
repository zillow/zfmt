package zfmtavro

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/heetch/avro/avrotypegen"
	"github.com/stretchr/testify/require"

	"github.com/zillow/zfmt/zfmtavro/testdata/heetch"
)

func demo() heetch.DemoSchema {
	return heetch.DemoSchema{
		IntField:    123,
		DoubleField: 123.4,
		StringField: "12345",
		BoolField:   true,
		BytesField:  []byte("123456"),
	}
}

// badSchemaRecord is a heetch record whose embedded schema is not valid avro.
type badSchemaRecord struct{}

func (badSchemaRecord) AvroRecord() avrotypegen.RecordInfo {
	return avrotypegen.RecordInfo{Schema: "not an avro schema"}
}

func TestFormatter_RoundTrip(t *testing.T) {
	// arrange
	fmtr := NewFormatter()
	got := heetch.DemoSchema{}

	// act
	b, err := fmtr.Marshal(demo())
	require.NoError(t, err)
	err = fmtr.Unmarshal(b, &got)

	// assert
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(demo(), got), "Unmarshal() mismatch (-want +got)")
}

func TestFormatter_Marshal_PointerEncodesSameBytesAsValue(t *testing.T) {
	// arrange
	fmtr := NewFormatter()
	rec := demo()
	want, err := fmtr.Marshal(rec)
	require.NoError(t, err)

	// act
	got, err := fmtr.Marshal(&rec)

	// assert
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestFormatter_Marshal_Errors(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr string
	}{
		{name: "non avro input", input: "what?", wantErr: "is not a heetch/avro generated type"},
		{name: "nil pointer to a record", input: (*heetch.DemoSchema)(nil), wantErr: "zfmtavro: cannot marshal a nil *heetch.DemoSchema"},
		{name: "record with invalid schema", input: badSchemaRecord{}, wantErr: "zfmtavro: encoding zfmtavro.badSchemaRecord: invalid character"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := NewFormatter()

			// act
			_, err := fmtr.Marshal(tt.input)

			// assert
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestFormatter_Unmarshal_Errors(t *testing.T) {
	valid, err := NewFormatter().Marshal(demo())
	require.NoError(t, err)

	tests := []struct {
		name    string
		data    []byte
		target  any
		wantErr string
	}{
		{name: "non avro target", data: valid, target: &bytes.Buffer{}, wantErr: "is not a heetch/avro generated type"},
		{name: "record with invalid schema", data: valid, target: &badSchemaRecord{}, wantErr: "zfmtavro: parsing schema of *zfmtavro.badSchemaRecord"},
		{name: "non pointer target", data: valid, target: heetch.DemoSchema{}, wantErr: "zfmtavro: cannot unmarshal into non-pointer heetch.DemoSchema"},
		{name: "typed nil target", data: valid, target: (*heetch.DemoSchema)(nil), wantErr: "zfmtavro: cannot unmarshal into a nil *heetch.DemoSchema"},
		{name: "nil payload", data: nil, target: &heetch.DemoSchema{}, wantErr: "zfmtavro: decoding into *heetch.DemoSchema"},
		{name: "truncated payload", data: valid[:3], target: &heetch.DemoSchema{}, wantErr: "zfmtavro: decoding into *heetch.DemoSchema"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := NewFormatter()

			// act
			err := fmtr.Unmarshal(tt.data, tt.target)

			// assert
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}
