package zfmtjson

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/zillow/zfmt/zfmtjson/testdata"
)

func TestFormatter_RoundTrip(t *testing.T) {
	// arrange
	fmtr := NewFormatter()
	in := testdata.ExampleJson{Id: "hello"}
	out := testdata.ExampleJson{}

	// act
	b, err := fmtr.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	err = fmtr.Unmarshal(b, &out)

	// assert
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"Id":"hello"}` {
		t.Errorf("unexpected json: %s", b)
	}
	if diff := cmp.Diff(in, out); diff != "" {
		t.Fatal(diff)
	}
}

func TestFormatter_Marshal_PropagatesEncodingError(t *testing.T) {
	// arrange
	fmtr := NewFormatter()

	// act
	_, err := fmtr.Marshal(make(chan int))

	// assert
	if err == nil {
		t.Fatal("expected error marshaling an unsupported type")
	}
}

func TestFormatter_Unmarshal_PropagatesDecodingError(t *testing.T) {
	// arrange
	fmtr := NewFormatter()
	out := testdata.ExampleJson{}

	// act
	err := fmtr.Unmarshal([]byte(`{"Id":`), &out)

	// assert
	if err == nil {
		t.Fatal("expected error unmarshaling invalid json")
	}
}
