package zfmt

import (
	"bytes"
	"reflect"
	"testing"
)

func TestStringFormatter_Marshall(t *testing.T) {
	type args struct {
		i any
	}
	tests := []struct {
		name    string
		f       *StringFormatter
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "string",
			f:       &StringFormatter{},
			args:    args{i: "test"},
			want:    []byte("test"),
			wantErr: false,
		},
		{
			name:    "byte array",
			f:       &StringFormatter{},
			args:    args{i: []byte("test")},
			want:    []byte("test"),
			wantErr: false,
		},
		{
			name:    "bytes buffer",
			f:       &StringFormatter{},
			args:    args{i: bytes.NewBufferString("test")},
			want:    []byte("test"),
			wantErr: false,
		},
		{
			name:    "invalid type",
			f:       &StringFormatter{},
			args:    args{i: 123},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &StringFormatter{}
			got, err := f.Marshall(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringFormatter.Marshall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringFormatter.Marshall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringFormatter_Unmarshal(t *testing.T) {
	type args struct {
		b []byte
		i any
	}
	tests := []struct {
		name    string
		f       *StringFormatter
		args    args
		wantErr bool
	}{
		{
			name: "string is immutable so that doesn't work",
			f:    &StringFormatter{},
			args: args{
				b: []byte("test"),
				i: func(str string) *string { return &str }(""),
			},
			wantErr: true,
		},
		{
			name: "supply io.Writer",
			f:    &StringFormatter{},
			args: args{
				b: []byte("test"),
				i: new(bytes.Buffer),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &StringFormatter{}
			if err := f.Unmarshal(tt.args.b, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("StringFormatter.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
