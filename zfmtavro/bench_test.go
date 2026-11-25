package zfmtavro

import (
	"testing"

	"github.com/zillow/zfmt/zfmtavro/testdata/heetch"
)

func BenchmarkFormatter_Marshal(b *testing.B) {
	fmtr := NewFormatter()
	rec := demo()
	for b.Loop() {
		if _, err := fmtr.Marshal(rec); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFormatter_Unmarshal(b *testing.B) {
	fmtr := NewFormatter()
	data, err := fmtr.Marshal(demo())
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		var got heetch.DemoSchema
		if err := fmtr.Unmarshal(data, &got); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSchematizedFormatter_Unmarshal(b *testing.B) {
	fmtr := NewSchematizedFormatter(99)
	data, err := fmtr.Marshal(demo())
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		var got heetch.DemoSchema
		if err := fmtr.Unmarshal(data, &got); err != nil {
			b.Fatal(err)
		}
	}
}
