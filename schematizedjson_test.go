package zfmt

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	av "gitlab.zgtools.net/devex/archetypes/gomods/zfmt/testdata"
)

func TestSchematizedJsonFormatter_Marshall(t *testing.T) {
	type fields struct {
		protobufFormatter JSONFormatter
		SchemaID          int
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "accept reference type of an object, default formatter should be useful",
			fields: fields{},
			args: args{
				v: &av.ExampleJson{
					Id: "hello",
				},
			},
			wantErr: false,
		},
		{
			name: "accept reference type of an object, with avro formatter and schemaID",
			fields: fields{
				protobufFormatter: JSONFormatter{},
				SchemaID:          99,
			},
			args: args{
				v: &av.ExampleJson{
					Id: "hello",
				},
			},
			wantErr: false,
		},
		{
			name: "accept value type of an object, with avro formatter and schemaID",
			fields: fields{
				protobufFormatter: JSONFormatter{},
				SchemaID:          99,
			},
			args: args{
				v: av.ExampleJson{
					Id: "hello",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SchematizedJSONFormatter{
				formatter: tt.fields.protobufFormatter,
				SchemaID:  tt.fields.SchemaID,
			}
			_, err := p.Marshall(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("SchematizedProtoFormatter.Marshall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSchematizedJsonFormmater_Unmarshall_PointerType(t *testing.T) {
	fmtter := SchematizedJSONFormatter{
		formatter: JSONFormatter{},
		SchemaID:  123,
	}
	expected := av.ExampleJson{Id: "what"}
	b, err := fmtter.Marshall(&expected)
	if err != nil {
		t.Fatal(err)
	}
	reslt := av.ExampleJson{}
	if err := fmtter.Unmarshal(b, &reslt); err != nil {
		t.Fatal(err)
	}
	diff := cmp.Diff(&reslt, &expected)
	if diff != "" {
		t.Fatal(diff)
	}
}
func TestSchematizedJsonFormmater_Unmarshalll_ValueType(t *testing.T) {
	fmtter := SchematizedJSONFormatter{
		formatter: JSONFormatter{},
		SchemaID:  123,
	}
	expected := av.ExampleJson{Id: "what"}
	b, err := fmtter.Marshall(expected)
	if err != nil {
		t.Fatal(err)
	}
	reslt := av.ExampleJson{}
	if err := fmtter.Unmarshal(b, &reslt); err != nil {
		t.Fatal(err)
	}
	diff := cmp.Diff(&reslt, &expected)
	if diff != "" {
		t.Fatal(diff)
	}
}
