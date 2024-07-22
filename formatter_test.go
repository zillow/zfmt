package zfmt

import (
	"reflect"
	"testing"
)

func TestGetFormatter(t *testing.T) {
	type args struct {
		fmt      FormatterType
		schemaID int
	}
	tests := []struct {
		name    string
		args    args
		want    Formatter
		wantErr bool
	}{
		{
			name: "json",
			args: args{
				fmt:      "json",
				schemaID: 0,
			},
			want: &JSONFormatter{},
		},
		{
			name: "proto_base64",
			args: args{
				fmt:      "proto_base64",
				schemaID: 0,
			},
			want: &ProtobufBase64Formatter{},
		},
		{
			name: "proto_raw",
			args: args{
				fmt:      "proto_raw",
				schemaID: 111,
			},
			want: &ProtobufRawFormatter{},
		},
		{
			name: "proto_json",
			args: args{
				fmt:      "proto_json",
				schemaID: 111,
			},
			want: &ProtoJSONFormatter{},
		},
		{
			name: "avro",
			args: args{
				fmt:      "avro",
				schemaID: 0,
			},
			want: &AvroFormatter{},
		},
		{
			name: "json_schema",
			args: args{
				fmt:      "json_schema",
				schemaID: 123,
			},
			want: &SchematizedJSONFormatter{SchemaID: 123},
		},
		{
			name: "avro_schema",
			args: args{
				fmt:      "avro_schema",
				schemaID: 789,
			},
			want: &SchematizedAvroFormatter{SchemaID: 789},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFormatter(tt.args.fmt, tt.args.schemaID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFormatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFormatter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
