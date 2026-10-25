package zfmtavro

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"github.com/zillow/zfmt/zfmtavro/testdata/heetch"
)

func TestSchematizedFormatter_Marshal_PrefixesWireFormatHeader(t *testing.T) {
	// arrange
	fmtr := NewSchematizedFormatter(0x01020304)
	payload, err := NewFormatter().Marshal(demo())
	require.NoError(t, err)

	// act
	got, err := fmtr.Marshal(demo())

	// assert
	require.NoError(t, err)
	require.Equal(t, append([]byte{0, 1, 2, 3, 4}, payload...), got)
}

func TestSchematizedFormatter_Marshal_PropagatesFormatterError(t *testing.T) {
	// arrange
	fmtr := NewSchematizedFormatter(99)

	// act
	_, err := fmtr.Marshal("what?")

	// assert
	require.ErrorContains(t, err, "zfmtavro: string is not a heetch/avro generated type")
}

func TestSchematizedFormatter_RoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		writer SchematizedFormatter
		reader SchematizedFormatter
	}{
		{name: "matching schema ID", writer: NewSchematizedFormatter(99), reader: NewSchematizedFormatter(99)},
		{name: "reader schema ID 0 skips validation", writer: NewSchematizedFormatter(99), reader: NewSchematizedFormatter(0)},
		{name: "schema ID 0 on both sides", writer: NewSchematizedFormatter(0), reader: NewSchematizedFormatter(0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			data, err := tt.writer.Marshal(demo())
			require.NoError(t, err)
			got := heetch.DemoSchema{}

			// act
			err = tt.reader.Unmarshal(data, &got)

			// assert
			require.NoError(t, err)
			require.Empty(t, cmp.Diff(demo(), got), "Unmarshal() mismatch (-want +got)")
		})
	}
}

func TestSchematizedFormatter_RoundTrip_PointerInput(t *testing.T) {
	// arrange
	fmtr := NewSchematizedFormatter(99)
	rec := demo()
	got := heetch.DemoSchema{}

	// act
	data, err := fmtr.Marshal(&rec)
	require.NoError(t, err)
	err = fmtr.Unmarshal(data, &got)

	// assert
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(demo(), got), "Unmarshal() mismatch (-want +got)")
}

func TestSchematizedFormatter_Unmarshal_Errors(t *testing.T) {
	valid, err := NewSchematizedFormatter(99).Marshal(demo())
	require.NoError(t, err)
	badMagic := append([]byte{1}, valid[1:]...)

	tests := []struct {
		name     string
		schemaID int
		data     []byte
		target   any
		wantErr  string
	}{
		{name: "payload shorter than the header", schemaID: 0, data: []byte{0, 1, 2}, target: &heetch.DemoSchema{}, wantErr: "shorter than the 5 byte wire format header"},
		{name: "non zero magic byte", schemaID: 0, data: badMagic, target: &heetch.DemoSchema{}, wantErr: "unknown magic byte 1"},
		{name: "schema ID mismatch", schemaID: 100, data: valid, target: &heetch.DemoSchema{}, wantErr: "schema ID mismatch, expected 100, got 99"},
		{name: "non avro target", schemaID: 99, data: valid, target: &[]byte{}, wantErr: "is not a heetch/avro generated type"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := NewSchematizedFormatter(tt.schemaID)

			// act
			err := fmtr.Unmarshal(tt.data, tt.target)

			// assert
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}
