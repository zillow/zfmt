package zfmtproto

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/zillow/zfmt/zfmtproto/testdata/example"
)

func TestSchematizedBase64FormatterDeprecated_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		fmtr    SchematizedBase64FormatterDeprecated
		input   any
		wantErr bool
	}{
		{
			name:    "zero value formatter accepts a proto pointer",
			fmtr:    SchematizedBase64FormatterDeprecated{},
			input:   &example.ExampleDef{Allowed: "!23", Disallowed: 7},
			wantErr: false,
		},
		{
			name:    "constructed formatter accepts a proto pointer",
			fmtr:    NewSchematizedBase64FormatterDeprecated(99),
			input:   &example.ExampleDef{Allowed: "!23", Disallowed: 7},
			wantErr: false,
		},
		{
			name:    "rejects a proto value (not a proto.Message)",
			fmtr:    NewSchematizedBase64FormatterDeprecated(99),
			input:   example.ExampleDef{Allowed: "!23", Disallowed: 7},
			wantErr: true,
		},
		{
			name:    "rejects a non proto type",
			fmtr:    NewSchematizedBase64FormatterDeprecated(99),
			input:   "what?",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := tt.fmtr

			// act
			_, err := fmtr.Marshal(tt.input)

			// assert
			require.Equal(t, tt.wantErr, err != nil, "unexpected error state: %v", err)
		})
	}
}

func TestSchematizedBase64FormatterDeprecated_RoundTrip(t *testing.T) {
	// arrange
	fmtr := NewSchematizedBase64FormatterDeprecated(123)
	in := &example.ExampleDef{Allowed: "1243", Disallowed: 987}
	out := &example.ExampleDef{}

	// act
	b, err := fmtr.Marshal(in)
	require.NoError(t, err)
	err = fmtr.Unmarshal(b, out)

	// assert
	require.NoError(t, err)
	require.True(t, proto.Equal(in, out), "round trip mismatch: %v != %v", in, out)
}

func TestSchematizedBase64FormatterDeprecated_Unmarshal_SchemaIDMismatch(t *testing.T) {
	// arrange
	writer := NewSchematizedBase64FormatterDeprecated(123)
	reader := NewSchematizedBase64FormatterDeprecated(124)
	b, err := writer.Marshal(&example.ExampleDef{Allowed: "1243"})
	require.NoError(t, err)

	// act
	err = reader.Unmarshal(b, &example.ExampleDef{})

	// assert
	require.Error(t, err)
}
