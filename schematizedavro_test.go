package zfmt

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/zillow/zfmt/testdata/heetch"

	av "github.com/zillow/zfmt/testdata"
)

func TestSchematizedAvroFormatter_Marshall(t *testing.T) {
	type fields struct {
		avroFmt  AvroFormatter
		SchemaID int
	}
	type args struct {
		v any
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "accept reference type of an avro object, default formatter should be useful",
			fields: fields{},
			args: args{
				v: &av.DemoSchema{
					IntField:    123,
					DoubleField: 123.4,
					StringField: "12345",
					BoolField:   true,
					BytesField:  []byte("123456"),
				},
			},
			wantErr: false,
		},
		{
			name: "accept reference type of an avro object, with avro formatter and schemaID",
			fields: fields{
				avroFmt:  AvroFormatter{},
				SchemaID: 99,
			},
			args: args{
				v: &av.DemoSchema{
					IntField:    123,
					DoubleField: 123.4,
					StringField: "12345",
					BoolField:   true,
					BytesField:  []byte("123456"),
				},
			},
			wantErr: false,
		},
		{
			name: "accept heetch avrorecord type",
			fields: fields{
				avroFmt:  AvroFormatter{},
				SchemaID: 99,
			},
			args: args{
				v: heetch.DemoSchema{
					IntField:    123,
					DoubleField: 123.4,
					StringField: "12345",
					BoolField:   true,
					BytesField:  []byte("123456"),
				},
			},
			wantErr: false,
		},
		{
			name: "accept value type of an avro object",
			fields: fields{
				avroFmt:  AvroFormatter{},
				SchemaID: 99,
			},
			args: args{
				v: av.DemoSchema{
					IntField:    123,
					DoubleField: 123.4,
					StringField: "12345",
					BoolField:   true,
					BytesField:  []byte("123456"),
				},
			},
			wantErr: false,
		},
		{
			name: "do not accept random type",
			fields: fields{
				avroFmt:  AvroFormatter{},
				SchemaID: 99,
			},
			args:    args{v: "what?"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SchematizedAvroFormatter{
				formatter: tt.fields.avroFmt,
				SchemaID:  tt.fields.SchemaID,
			}
			_, err := p.Marshall(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("SchematizedAvroFormatter.Marshall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSchematizedAvroFormatter_UnmarshalToInvalidAvroType(t *testing.T) {
	fmtter := &SchematizedAvroFormatter{SchemaID: 99}
	input := &av.DemoSchema{
		IntField:    123,
		DoubleField: 123.4,
		StringField: "12345",
		BoolField:   true,
		BytesField:  []byte("123456"),
	}
	data, err := fmtter.Marshall(input)
	if err != nil {
		t.Errorf("should not have error marshalling %v", err)
	}
	var output bytes.Buffer
	err = fmtter.Unmarshal(data, &output)
	if err == nil {
		t.Errorf("should get error because output type is not an avro object")
	}
}

func TestSchematizedAvroFormatter_UnmarshalNonSchematizedAvro(t *testing.T) {
	binAvro := []byte{}
	var output bytes.Buffer
	outFmtter := &SchematizedAvroFormatter{}
	err := outFmtter.Unmarshal(binAvro, &output)
	if err == nil {
		t.Errorf("should get error because input does not contain schema")
	}
}

func TestSchematizedAvroFormatter_UnmarshalValidAvroWithSchemaID(t *testing.T) {
	type testCase struct {
		Name     string
		Input    any
		Expected any
		Output   any
	}

	testCases := []testCase{
		{
			Name: "gogen-avro with with schemaID",
			Input: &av.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Expected: &av.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Output: &av.DemoSchema{},
		},
		{
			Name: "heetch avro with with schemaID",
			Input: heetch.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Expected: &heetch.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Output: &heetch.DemoSchema{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fmtter := &SchematizedAvroFormatter{SchemaID: 99}
			input := tc.Input
			data, err := fmtter.Marshall(input)
			if err != nil {
				t.Errorf("should not have error marshalling %v", err)
			}
			output := tc.Output
			err = fmtter.Unmarshal(data, output)
			if err != nil {
				t.Errorf("should not have error marshalling %v", err)
			}
			if !reflect.DeepEqual(tc.Expected, output) {
				t.Errorf("data should match, want=%v, got=%v", input, output)
			}
		})
	}
}

func TestSchematizedAvroFormatter_UnmarshalValidAvroWithNoSchemaID(t *testing.T) {
	type testCase struct {
		Name     string
		Input    any
		Expected any
		Output   any
	}

	testCases := []testCase{
		{
			Name: "gogen-avro with no schemaID",
			Input: &av.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Expected: &av.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Output: &av.DemoSchema{},
		},
		{
			Name: "heetch avro with no schemaID",
			Input: heetch.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Expected: &heetch.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Output: &heetch.DemoSchema{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			inFmtter := &SchematizedAvroFormatter{SchemaID: 99}
			input := tc.Input
			data, err := inFmtter.Marshall(input)
			if err != nil {
				t.Errorf("should not have error marshalling %v", err)
			}

			// when schemaID is unset, as long as the data is unmarshallable, we don't throw error
			outFmtter := &SchematizedAvroFormatter{}
			output := tc.Output
			err = outFmtter.Unmarshal(data, output)
			if err != nil {
				t.Errorf("should not have error marshalling %v", err)
			}
			if !reflect.DeepEqual(tc.Expected, output) {
				t.Errorf("data should match, want=%v, got=%v", input, output)
			}
		})
	}

}

func TestSchematizedAvroFormatter_UnmarshalValidAvroWithWrongSchemaID(t *testing.T) {
	type testCase struct {
		Name   string
		Input  any
		Output any
	}

	testCases := []testCase{
		{
			Name: "gogen-avro with wrong schemaID",
			Input: &av.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Output: &av.DemoSchema{},
		},
		{
			Name: "heetch avro with wrong schemaID",
			Input: heetch.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Output: &heetch.DemoSchema{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			inFmtter := &SchematizedAvroFormatter{SchemaID: 99}
			data, err := inFmtter.Marshall(tc.Input)
			if err != nil {
				t.Errorf("should not have error marshalling %v", err)
			}

			outFmtter := &SchematizedAvroFormatter{SchemaID: 100}
			err = outFmtter.Unmarshal(data, tc.Output)
			if err == nil {
				t.Error("should have error unmarshalling due to incorrect schema ID")
			}
		})
	}
}

func TestSchematizedAvroFormatter_Equivalency(t *testing.T) {
	type testCase struct {
		Name     string
		Input    any
		Expected any
		Output   any
	}

	testCases := []testCase{
		{
			Name: "marshal gogen-avro unmarshal heetch avro",
			Input: &av.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Expected: &heetch.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Output: &heetch.DemoSchema{},
		},
		{
			Name: "marshall heetch unmarshall gogen-avro",
			Input: heetch.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Expected: &av.DemoSchema{
				IntField:    123,
				DoubleField: 123.4,
				StringField: "12345",
				BoolField:   true,
				BytesField:  []byte("123456"),
			},
			Output: &av.DemoSchema{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fmtter := &SchematizedAvroFormatter{SchemaID: 99}
			input := tc.Input
			data, err := fmtter.Marshall(input)
			if err != nil {
				t.Errorf("should not have error marshalling %v", err)
			}
			output := tc.Output
			err = fmtter.Unmarshal(data, output)
			if err != nil {
				t.Errorf("should not have error marshalling %v", err)
			}
			if !reflect.DeepEqual(tc.Expected, output) {
				t.Errorf("data should match, want=%v, got=%v", input, output)
			}
		})
	}
}
