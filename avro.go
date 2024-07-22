package zfmt

import (
	"bytes"
	"fmt"

	"github.com/actgardner/gogen-avro/v10/compiler"
	"github.com/actgardner/gogen-avro/v10/soe"
	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/actgardner/gogen-avro/v10/vm/types"
	heetch "github.com/heetch/avro"
	heetchtypegen "github.com/heetch/avro/avrotypegen"
)

// GeneratedAvroRecord combines interfaces that make Encoding/Decoding possible
// for gogen-avro struct
type GeneratedAvroRecord interface {
	soe.AvroRecord
	types.Field
	Schema() string
}

type AvroFormatter struct{}

func (p *AvroFormatter) Marshall(v any) ([]byte, error) {
	switch m := v.(type) {
	case soe.AvroRecord:
		buf := &bytes.Buffer{}
		err := m.Serialize(buf)
		return buf.Bytes(), err
	case heetchtypegen.AvroRecord:
		b, _, err := heetch.Marshal(v)
		return b, err
	default:
		return nil, fmt.Errorf("%T, avro formatter supports only gogen-avro or heetch avro messages", v)
	}
}

func (p *AvroFormatter) Unmarshal(b []byte, v any) error {
	switch m := v.(type) {
	case GeneratedAvroRecord:
		r := bytes.NewReader(b)
		deser, err := compiler.CompileSchemaBytes([]byte(m.Schema()), []byte(m.Schema()))
		if err != nil {
			return err
		}
		return vm.Eval(r, deser, m)
	case heetchtypegen.AvroRecord:
		t, err := heetch.ParseType(m.AvroRecord().Schema)
		if err != nil {
			return err
		}
		_, err = heetch.Unmarshal(b, v, t)
		return err
	default:
		return fmt.Errorf("%T, avro formatter supports only gogen-avro or heetch avro messages", v)
	}
}
