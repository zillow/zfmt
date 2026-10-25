package zfmtstring

import (
	"bytes"
	"reflect"
	"testing"
)

func TestStringFormatter_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    []byte
		wantErr bool
	}{
		{name: "string", input: "test", want: []byte("test")},
		{name: "byte slice", input: []byte("test"), want: []byte("test")},
		{name: "io.Reader", input: bytes.NewBufferString("test"), want: []byte("test")},
		{name: "unsupported type", input: 123, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := NewStringFormatter()

			// act
			got, err := fmtr.Marshal(tt.input)

			// assert
			if (err != nil) != tt.wantErr {
				t.Fatalf("StringFormatter.Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringFormatter.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringFormatter_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		target  any
		wantErr bool
	}{
		{name: "string pointer is rejected because strings are immutable", target: new(string), wantErr: true},
		{name: "io.Writer is accepted", target: new(bytes.Buffer)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fmtr := NewStringFormatter()

			// act
			err := fmtr.Unmarshal([]byte("test"), tt.target)

			// assert
			if (err != nil) != tt.wantErr {
				t.Fatalf("StringFormatter.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if w, ok := tt.target.(*bytes.Buffer); ok && w.String() != "test" {
				t.Errorf("writer received %q, want %q", w.String(), "test")
			}
		})
	}
}
