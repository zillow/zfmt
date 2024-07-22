package zfmt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	example2 "gitlab.zgtools.net/devex/archetypes/gomods/zfmt/testdata/example"
)

func TestSchematizedProtoDeprecatedFormatter_Marshall(t *testing.T) {
	type fields struct {
		protobufFormatter ProtobufBase64Formatter
		SchemaID          int
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
			name:   "accept reference type of an proto object, default formatter should be useful",
			fields: fields{},
			args: args{
				v: &example2.ExampleDef{
					Allowed:    "!23",
					Disallowed: 7,
				},
			},
			wantErr: false,
		},
		{
			name: "accept reference type of an proto object, with avro formatter and schemaID",
			fields: fields{
				protobufFormatter: ProtobufBase64Formatter{},
				SchemaID:          99,
			},
			args: args{
				v: &example2.ExampleDef{
					Allowed:    "!23",
					Disallowed: 7,
				},
			},
			wantErr: false,
		},
		{
			name: "do not accept value type of a proto object",
			fields: fields{
				protobufFormatter: ProtobufBase64Formatter{},
				SchemaID:          99,
			},
			args: args{
				v: example2.ExampleDef{
					Allowed:    "!23",
					Disallowed: 7,
				},
			},
			wantErr: true,
		},
		{
			name: "do not accept random type",
			fields: fields{
				protobufFormatter: ProtobufBase64Formatter{},
				SchemaID:          99,
			},
			args:    args{v: "what?"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SchematizedProtoFormatterDeprecated{
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

func TestSchematizedProtoDeprecatedFormmater_Unmarshall(t *testing.T) {
	fmtter := SchematizedProtoFormatterDeprecated{
		formatter: ProtobufBase64Formatter{},
		SchemaID:  123,
	}
	expected := example2.ExampleDef{Allowed: "1243", Disallowed: 987}
	b, err := fmtter.Marshall(&expected)
	if err != nil {
		t.Fatal(err)
	}
	reslt := example2.ExampleDef{}
	if err := fmtter.Unmarshal(b, &reslt); err != nil {
		t.Fatal(err)
	}
	diff := cmp.Diff(&reslt, &expected, cmpopts.IgnoreUnexported(example2.ExampleDef{}))
	if diff != "" {
		t.Fatal(diff)
	}
}
