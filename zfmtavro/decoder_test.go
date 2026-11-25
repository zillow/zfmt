package zfmtavro

import (
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/heetch/avro/avrotypegen"
	"github.com/stretchr/testify/require"

	"github.com/zillow/zfmt/zfmtavro/testdata/heetch"
)

// otherRecord is a second record type whose schema differs from heetch.DemoSchema's.
type otherRecord struct {
	Name string
}

func (otherRecord) AvroRecord() avrotypegen.RecordInfo {
	return avrotypegen.RecordInfo{
		Schema:   `{"type":"record","name":"Other","fields":[{"name":"Name","type":"string"}]}`,
		Required: []bool{0: true},
	}
}

type unmarshaler interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(b []byte, v any) error
}

func TestFormatters_DecodeDifferentRecordTypesWithTheirOwnSchemas(t *testing.T) {
	// A heetch SingleDecoder caches writer schemas by schema ID alone. Both record types below are decoded
	// through one formatter with the same schema ID, which fails if the decoders share that cache.
	tests := []struct {
		name string
		fmtr unmarshaler
	}{
		{name: "Formatter", fmtr: NewFormatter()},
		{name: "SchematizedFormatter", fmtr: NewSchematizedFormatter(0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			demoData, err := tt.fmtr.Marshal(demo())
			require.NoError(t, err)
			otherData, err := tt.fmtr.Marshal(otherRecord{Name: "other"})
			require.NoError(t, err)
			var gotDemo heetch.DemoSchema
			var gotOther otherRecord

			// act
			errDemo := tt.fmtr.Unmarshal(demoData, &gotDemo)
			errOther := tt.fmtr.Unmarshal(otherData, &gotOther)

			// assert
			require.NoError(t, errDemo)
			require.NoError(t, errOther)
			require.Empty(t, cmp.Diff(demo(), gotDemo), "Unmarshal() mismatch (-want +got)")
			require.Equal(t, otherRecord{Name: "other"}, gotOther)
		})
	}
}

func TestDecoders_ReuseDecoderPerType(t *testing.T) {
	// arrange
	fmtr := NewFormatter()
	data, err := fmtr.Marshal(demo())
	require.NoError(t, err)
	var first, second heetch.DemoSchema

	// act
	require.NoError(t, fmtr.Unmarshal(data, &first))
	require.NoError(t, fmtr.Unmarshal(data, &second))

	// assert
	fmtr.decoders.mu.RLock()
	defer fmtr.decoders.mu.RUnlock()
	require.Len(t, fmtr.decoders.byType, 1)
}

func TestDecoders_CachesAreNotSharedBetweenFormatters(t *testing.T) {
	// arrange
	a, b := NewFormatter(), NewFormatter()
	data, err := a.Marshal(demo())
	require.NoError(t, err)
	var got heetch.DemoSchema

	// act
	err = a.Unmarshal(data, &got)

	// assert
	require.NoError(t, err)
	require.Len(t, a.decoders.byType, 1)
	require.Empty(t, b.decoders.byType)
}

func TestDecoders_DoNotCacheInvalidSchema(t *testing.T) {
	// arrange
	fmtr := NewFormatter()

	// act
	err := fmtr.Unmarshal([]byte{0}, &badSchemaRecord{})

	// assert
	require.ErrorContains(t, err, "parsing schema of *zfmtavro.badSchemaRecord")
	require.Empty(t, fmtr.decoders.byType)
}

func TestFormatter_Unmarshal_ConcurrentDecodes(t *testing.T) {
	// arrange
	fmtr := NewFormatter()
	data, err := fmtr.Marshal(demo())
	require.NoError(t, err)
	const workers = 16
	errs := make([]error, workers)
	var wg sync.WaitGroup

	// act
	for i := range workers {
		wg.Go(func() {
			var got heetch.DemoSchema
			errs[i] = fmtr.Unmarshal(data, &got)
		})
	}
	wg.Wait()

	// assert
	for _, err := range errs {
		require.NoError(t, err)
	}
}

func TestZeroValueFormatters_Unmarshal_ReturnError(t *testing.T) {
	tests := []struct {
		name string
		fmtr unmarshaler
	}{
		{name: "Formatter", fmtr: Formatter{}},
		{name: "SchematizedFormatter", fmtr: SchematizedFormatter{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			data := []byte{0, 0, 0, 0, 0}

			// act
			err := tt.fmtr.Unmarshal(data, &heetch.DemoSchema{})

			// assert
			require.ErrorIs(t, err, errNotConstructed)
		})
	}
}
