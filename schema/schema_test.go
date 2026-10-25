package schema

import (
	"bytes"
	"encoding/json"
	"testing"
)

type jsonFormatter struct{}

func (jsonFormatter) Marshal(v any) ([]byte, error)   { return json.Marshal(v) }
func (jsonFormatter) Unmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }

type payload struct {
	ID string `json:"id"`
}

func TestMarshal_PrependsConfluentHeader(t *testing.T) {
	// arrange
	wantHeader := []byte{0, 1, 2, 3, 4}

	// act
	got, err := Marshal(jsonFormatter{}, 0x01020304, payload{ID: "x"})

	// assert
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got[:5], wantHeader) {
		t.Fatalf("header = %v, want %v", got[:5], wantHeader)
	}
	if string(got[5:]) != `{"id":"x"}` {
		t.Fatalf("body = %q", got[5:])
	}
}

func TestUnmarshal_RoundTrip(t *testing.T) {
	// arrange
	b, err := Marshal(jsonFormatter{}, 99, payload{ID: "x"})
	if err != nil {
		t.Fatal(err)
	}
	var out payload

	// act
	err = Unmarshal(jsonFormatter{}, 99, b, &out)

	// assert
	if err != nil {
		t.Fatal(err)
	}
	if out.ID != "x" {
		t.Fatalf("out = %+v", out)
	}
}

func TestUnmarshal_ZeroSchemaIDSkipsValidation(t *testing.T) {
	// arrange
	b, err := Marshal(jsonFormatter{}, 99, payload{ID: "x"})
	if err != nil {
		t.Fatal(err)
	}
	var out payload

	// act
	err = Unmarshal(jsonFormatter{}, 0, b, &out)

	// assert
	if err != nil {
		t.Fatalf("expected zero schema ID to bypass validation, got %v", err)
	}
}

func TestUnmarshal_SchemaIDMismatch(t *testing.T) {
	// arrange
	b, err := Marshal(jsonFormatter{}, 99, payload{ID: "x"})
	if err != nil {
		t.Fatal(err)
	}
	var out payload

	// act
	err = Unmarshal(jsonFormatter{}, 100, b, &out)

	// assert
	if err == nil {
		t.Fatal("expected schema ID mismatch error")
	}
}

func TestUnmarshal_TooShort(t *testing.T) {
	// arrange
	var out payload

	// act
	err := Unmarshal(jsonFormatter{}, 0, []byte{0, 1, 2}, &out)

	// assert
	if err == nil {
		t.Fatal("expected error for payload shorter than the wire format header")
	}
}

func TestMarshal_PropagatesFormatterError(t *testing.T) {
	// arrange
	unencodable := make(chan int)

	// act
	_, err := Marshal(jsonFormatter{}, 1, unencodable)

	// assert
	if err == nil {
		t.Fatal("expected json marshal error to propagate")
	}
}
