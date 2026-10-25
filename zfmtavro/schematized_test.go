package zfmtavro

import (
	"bytes"
	"reflect"
	"testing"

	av "github.com/zillow/zfmt/zfmtavro/testdata"
	"github.com/zillow/zfmt/zfmtavro/testdata/heetch"
)

func TestSchematizedFormatter_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		fmtr    SchematizedFormatter
		input   any
		wantErr bool
	}{
		{
			name:    "zero value formatter accepts a gogen-avro pointer",
			fmtr:    SchematizedFormatter{},
			input:   gogenDemo(),
			wantErr: false,
		},
		{
			name:    "constructed formatter accepts a gogen-avro pointer",
			fmtr:    NewSchematizedFormatter(99),
			input:   gogenDemo(),
			wantErr: false,
		},
		{
			name:    "constructed formatter accepts a gogen-avro value",
			fmtr:    NewSchematizedFormatter(99),
			input:   *gogenDemo(),
			wantErr: false,
		},
		{
			name:    "constructed formatter accepts a heetch record",
			fmtr:    NewSchematizedFormatter(99),
			input:   heetchDemo(),
			wantErr: false,
		},
		{
			name:    "rejects a non avro type",
			fmtr:    NewSchematizedFormatter(99),
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
			if (err != nil) != tt.wantErr {
				t.Errorf("SchematizedFormatter.Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSchematizedFormatter_Unmarshal_RejectsNonAvroTarget(t *testing.T) {
	// arrange
	fmtr := NewSchematizedFormatter(99)
	data, err := fmtr.Marshal(gogenDemo())
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer

	// act
	err = fmtr.Unmarshal(data, &output)

	// assert
	if err == nil {
		t.Error("should get error because output type is not an avro object")
	}
}

func TestSchematizedFormatter_Unmarshal_RejectsPayloadWithoutSchemaHeader(t *testing.T) {
	// arrange
	fmtr := SchematizedFormatter{}
	var output bytes.Buffer

	// act
	err := fmtr.Unmarshal([]byte{}, &output)

	// assert
	if err == nil {
		t.Error("should get error because input does not contain schema")
	}
}

func TestSchematizedFormatter_RoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		writer   SchematizedFormatter
		reader   SchematizedFormatter
		input    any
		expected any
		output   any
	}{
		{
			name:     "gogen-avro, matching schema ID",
			writer:   NewSchematizedFormatter(99),
			reader:   NewSchematizedFormatter(99),
			input:    gogenDemo(),
			expected: gogenDemo(),
			output:   &av.DemoSchema{},
		},
		{
			name:     "heetch avro, matching schema ID",
			writer:   NewSchematizedFormatter(99),
			reader:   NewSchematizedFormatter(99),
			input:    heetchDemo(),
			expected: ptr(heetchDemo()),
			output:   &heetch.DemoSchema{},
		},
		{
			name:     "gogen-avro, reader schema ID unset skips validation",
			writer:   NewSchematizedFormatter(99),
			reader:   SchematizedFormatter{},
			input:    gogenDemo(),
			expected: gogenDemo(),
			output:   &av.DemoSchema{},
		},
		{
			name:     "heetch avro, reader schema ID unset skips validation",
			writer:   NewSchematizedFormatter(99),
			reader:   SchematizedFormatter{},
			input:    heetchDemo(),
			expected: ptr(heetchDemo()),
			output:   &heetch.DemoSchema{},
		},
		{
			name:     "gogen-avro written, heetch read",
			writer:   NewSchematizedFormatter(99),
			reader:   NewSchematizedFormatter(99),
			input:    gogenDemo(),
			expected: ptr(heetchDemo()),
			output:   &heetch.DemoSchema{},
		},
		{
			name:     "heetch written, gogen-avro read",
			writer:   NewSchematizedFormatter(99),
			reader:   NewSchematizedFormatter(99),
			input:    heetchDemo(),
			expected: gogenDemo(),
			output:   &av.DemoSchema{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			data, err := tt.writer.Marshal(tt.input)
			if err != nil {
				t.Fatalf("should not have error marshaling %v", err)
			}

			// act
			err = tt.reader.Unmarshal(data, tt.output)

			// assert
			if err != nil {
				t.Fatalf("should not have error unmarshaling %v", err)
			}
			if !reflect.DeepEqual(tt.expected, tt.output) {
				t.Errorf("data should match, want=%v, got=%v", tt.expected, tt.output)
			}
		})
	}
}

func TestSchematizedFormatter_Unmarshal_SchemaIDMismatch(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		output any
	}{
		{name: "gogen-avro", input: gogenDemo(), output: &av.DemoSchema{}},
		{name: "heetch avro", input: heetchDemo(), output: &heetch.DemoSchema{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			writer := NewSchematizedFormatter(99)
			reader := NewSchematizedFormatter(100)
			data, err := writer.Marshal(tt.input)
			if err != nil {
				t.Fatalf("should not have error marshaling %v", err)
			}

			// act
			err = reader.Unmarshal(data, tt.output)

			// assert
			if err == nil {
				t.Error("should have error unmarshaling due to incorrect schema ID")
			}
		})
	}
}
