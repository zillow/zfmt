package zfmtjson

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/zillow/zfmt/zfmtjson/testdata"
)

func TestSchematizedFormatter_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		fmtr    SchematizedFormatter
		input   any
		wantErr bool
	}{
		{
			name:    "zero value formatter accepts a pointer",
			fmtr:    SchematizedFormatter{},
			input:   &testdata.ExampleJson{Id: "hello"},
			wantErr: false,
		},
		{
			name:    "constructed formatter accepts a pointer",
			fmtr:    NewSchematizedFormatter(99),
			input:   &testdata.ExampleJson{Id: "hello"},
			wantErr: false,
		},
		{
			name:    "constructed formatter accepts a value",
			fmtr:    NewSchematizedFormatter(99),
			input:   testdata.ExampleJson{Id: "hello"},
			wantErr: false,
		},
		{
			name:    "propagates json encoding errors",
			fmtr:    NewSchematizedFormatter(99),
			input:   make(chan int),
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

func TestSchematizedFormatter_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{name: "pointer input", input: &testdata.ExampleJson{Id: "what"}},
		{name: "value input", input: testdata.ExampleJson{Id: "what"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := NewSchematizedFormatter(123)
			expected := testdata.ExampleJson{Id: "what"}
			out := testdata.ExampleJson{}

			// act
			b, err := fmtr.Marshal(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			err = fmtr.Unmarshal(b, &out)

			// assert
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(expected, out); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSchematizedFormatter_Unmarshal_SchemaIDMismatch(t *testing.T) {
	// arrange
	writer := NewSchematizedFormatter(123)
	reader := NewSchematizedFormatter(124)
	b, err := writer.Marshal(testdata.ExampleJson{Id: "what"})
	if err != nil {
		t.Fatal(err)
	}

	// act
	err = reader.Unmarshal(b, &testdata.ExampleJson{})

	// assert
	if err == nil {
		t.Fatal("expected schema ID mismatch error")
	}
}

func TestSchematizedFormatter_Unmarshal_ZeroSchemaIDSkipsValidation(t *testing.T) {
	// arrange
	writer := NewSchematizedFormatter(123)
	reader := SchematizedFormatter{}
	b, err := writer.Marshal(testdata.ExampleJson{Id: "what"})
	if err != nil {
		t.Fatal(err)
	}
	out := testdata.ExampleJson{}

	// act
	err = reader.Unmarshal(b, &out)

	// assert
	if err != nil {
		t.Fatal(err)
	}
	if out.Id != "what" {
		t.Fatalf("out = %+v", out)
	}
}
