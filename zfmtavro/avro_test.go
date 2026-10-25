package zfmtavro

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/heetch/avro/avrotypegen"

	av "github.com/zillow/zfmt/zfmtavro/testdata"
	"github.com/zillow/zfmt/zfmtavro/testdata/heetch"
)

func gogenDemo() *av.DemoSchema {
	return &av.DemoSchema{
		IntField:    123,
		DoubleField: 123.4,
		StringField: "12345",
		BoolField:   true,
		BytesField:  []byte("123456"),
	}
}

func heetchDemo() heetch.DemoSchema {
	return heetch.DemoSchema{
		IntField:    123,
		DoubleField: 123.4,
		StringField: "12345",
		BoolField:   true,
		BytesField:  []byte("123456"),
	}
}

// badSchemaGogen is a gogen-avro record whose embedded schema is not valid avro,
// so schema compilation fails during Unmarshal.
type badSchemaGogen struct {
	*av.DemoSchema
}

func (badSchemaGogen) Schema() string { return "not an avro schema" }

// badSchemaHeetch is a heetch record whose embedded schema is not valid avro.
type badSchemaHeetch struct{}

func (badSchemaHeetch) AvroRecord() avrotypegen.RecordInfo {
	return avrotypegen.RecordInfo{Schema: "not an avro schema"}
}

func TestFormatter_RoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
		output   any
	}{
		{name: "gogen-avro", input: gogenDemo(), expected: gogenDemo(), output: &av.DemoSchema{}},
		{name: "heetch avro", input: heetchDemo(), expected: ptr(heetchDemo()), output: &heetch.DemoSchema{}},
		{name: "gogen-avro written, heetch read", input: gogenDemo(), expected: ptr(heetchDemo()), output: &heetch.DemoSchema{}},
		{name: "heetch written, gogen-avro read", input: heetchDemo(), expected: gogenDemo(), output: &av.DemoSchema{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := NewFormatter()

			// act
			b, err := fmtr.Marshal(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			err = fmtr.Unmarshal(b, tt.output)

			// assert
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.expected, tt.output) {
				t.Errorf("data should match, want=%v, got=%v", tt.expected, tt.output)
			}
		})
	}
}

func TestFormatter_Marshal_RejectsNonAvro(t *testing.T) {
	// arrange
	fmtr := NewFormatter()

	// act
	_, err := fmtr.Marshal("what?")

	// assert
	if err == nil {
		t.Fatal("expected error for non avro input")
	}
}

func TestFormatter_Unmarshal_Errors(t *testing.T) {
	tests := []struct {
		name   string
		target any
	}{
		{name: "non avro target", target: &bytes.Buffer{}},
		{name: "gogen-avro record with invalid schema", target: badSchemaGogen{DemoSchema: &av.DemoSchema{}}},
		{name: "heetch record with invalid schema", target: badSchemaHeetch{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := NewFormatter()
			b, err := fmtr.Marshal(gogenDemo())
			if err != nil {
				t.Fatal(err)
			}

			// act
			err = fmtr.Unmarshal(b, tt.target)

			// assert
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func ptr[T any](v T) *T { return &v }
