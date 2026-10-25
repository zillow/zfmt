package zfmtavro

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/heetch/avro"
	"github.com/heetch/avro/avrotypegen"
)

var avroRecordType = reflect.TypeFor[avrotypegen.AvroRecord]()

// Formatter encodes/decodes heetch/avro generated types using the avro binary encoding.
//
// Create it with NewFormatter. It caches a compiled decoder per record type, and copies of a Formatter share
// that cache. Unmarshal assumes the message was written with the target type's own schema: messages written
// with a different, even compatible, schema are not resolved against it. Use a schema registry aware
// marshaler, such as zkafka's, when producers and consumers evolve schemas independently.
type Formatter struct {
	decoders *decoders
}

// NewFormatter returns an avro Formatter
func NewFormatter() Formatter {
	return Formatter{decoders: newDecoders(false)}
}

// Marshal encodes v, which must be a heetch/avro generated type or a pointer to one, as avro binary.
// heetch caches the compiled encoder per Go type, so repeated calls do not recompile it.
func (Formatter) Marshal(v any) ([]byte, error) {
	if _, ok := v.(avrotypegen.AvroRecord); !ok {
		return nil, fmt.Errorf("zfmtavro: %T is not a heetch/avro generated type", v)
	}
	// heetch reads the schema by calling AvroRecord on the zero value of v's type. For a pointer to a
	// generated record that zero value is a nil pointer, and the generated value-receiver method panics
	// on it, so marshal the record the pointer refers to instead. The encoded bytes are the same.
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Pointer && rv.Type().Elem().Implements(avroRecordType) {
		if rv.IsNil() {
			return nil, fmt.Errorf("zfmtavro: cannot marshal a nil %T", v)
		}
		v = rv.Elem().Interface()
	}
	b, _, err := avro.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("zfmtavro: encoding %T: %w", v, err)
	}
	return b, nil
}

// Unmarshal decodes avro binary into v, which must be a non-nil pointer to a heetch/avro generated type
func (f Formatter) Unmarshal(b []byte, v any) error {
	if f.decoders == nil {
		return errNotConstructed
	}
	return f.decoders.unmarshal(b, v)
}

var errNotConstructed = errors.New("zfmtavro: formatter must be created with its New... constructor")
