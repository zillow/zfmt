package zfmtavro

import (
	"context"
	"encoding/binary"
	"fmt"
	"reflect"
	"sync"

	"github.com/heetch/avro"
	"github.com/heetch/avro/avrotypegen"
)

// decoders holds one heetch SingleDecoder per Go record type. A SingleDecoder caches the compiled decode
// program, which is the dominant cost of decoding; heetch's package-level Unmarshal recompiles it on every
// call. Each type gets its own decoder because a SingleDecoder caches writer schemas by schema ID alone,
// while zfmtavro has no schema registry and answers every lookup with the target type's own schema.
type decoders struct {
	// framed reports whether messages start with the Confluent wire format header.
	framed bool

	mu     sync.RWMutex
	byType map[reflect.Type]*avro.SingleDecoder
}

func newDecoders(framed bool) *decoders {
	return &decoders{framed: framed, byType: map[reflect.Type]*avro.SingleDecoder{}}
}

// unmarshal decodes data into v, which must be a non-nil pointer to a heetch/avro generated type.
func (d *decoders) unmarshal(data []byte, v any) error {
	m, ok := v.(avrotypegen.AvroRecord)
	if !ok {
		return fmt.Errorf("zfmtavro: %T is not a heetch/avro generated type", v)
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("zfmtavro: cannot unmarshal into non-pointer %T", v)
	}
	if rv.IsNil() {
		return fmt.Errorf("zfmtavro: cannot unmarshal into a nil %T", v)
	}
	dec, err := d.decoderFor(rv.Type(), m)
	if err != nil {
		return fmt.Errorf("zfmtavro: parsing schema of %T: %w", v, err)
	}
	if data == nil {
		// heetch treats a nil message as one with no schema ID; decode it as an empty payload instead.
		data = []byte{}
	}
	if _, err := dec.Unmarshal(context.Background(), data, v); err != nil {
		return fmt.Errorf("zfmtavro: decoding into %T: %w", v, err)
	}
	return nil
}

func (d *decoders) decoderFor(t reflect.Type, m avrotypegen.AvroRecord) (*avro.SingleDecoder, error) {
	d.mu.RLock()
	dec, ok := d.byType[t]
	d.mu.RUnlock()
	if ok {
		return dec, nil
	}
	schema, err := avro.ParseType(m.AvroRecord().Schema)
	if err != nil {
		return nil, err
	}
	dec = avro.NewSingleDecoder(readerSchema{schema: schema, framed: d.framed}, nil)
	// Concurrent first decodes of the same type may each build a decoder; they are equivalent, so the last store wins.
	d.mu.Lock()
	d.byType[t] = dec
	d.mu.Unlock()
	return dec, nil
}

// readerSchema is a heetch DecodingRegistry that reports the target type's own schema as the writer
// schema of every message.
type readerSchema struct {
	schema *avro.Type
	framed bool
}

// DecodeSchemaID returns the schema ID from the Confluent header of a framed message, which the caller
// has already validated, or 0 for an unframed message.
func (r readerSchema) DecodeSchemaID(msg []byte) (int64, []byte) {
	if !r.framed {
		return 0, msg
	}
	return int64(binary.BigEndian.Uint32(msg[1:headerLen])), msg[headerLen:]
}

// SchemaForID returns the target type's schema for every ID.
func (r readerSchema) SchemaForID(context.Context, int64) (*avro.Type, error) {
	return r.schema, nil
}
