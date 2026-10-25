package zfmtbytes

import (
	"fmt"
	"io"
)

// Formatter passes payloads through as raw bytes. Marshal accepts a string, []byte or io.Reader.
// Unmarshal writes into a *string, *[]byte or io.Writer. The zero value is ready to use.
type Formatter struct{}

// Marshal returns the bytes of a string, []byte or io.Reader
func (Formatter) Marshal(v any) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	case io.Reader:
		b, err := io.ReadAll(v)
		if err != nil {
			return nil, fmt.Errorf("zfmtbytes: reading %T: %w", v, err)
		}
		return b, nil
	default:
		return nil, fmt.Errorf("zfmtbytes: cannot marshal %T, want string, []byte or io.Reader", v)
	}
}

// Unmarshal stores b in v, which must be a *string, *[]byte or io.Writer.
// A *[]byte target receives a copy, because transports may reuse the backing array of b.
func (Formatter) Unmarshal(b []byte, v any) error {
	switch v := v.(type) {
	case *string:
		if v == nil {
			return fmt.Errorf("zfmtbytes: cannot unmarshal into a nil %T", v)
		}
		*v = string(b)
		return nil
	case *[]byte:
		if v == nil {
			return fmt.Errorf("zfmtbytes: cannot unmarshal into a nil %T", v)
		}
		*v = append((*v)[:0], b...)
		return nil
	case io.Writer:
		if _, err := v.Write(b); err != nil {
			return fmt.Errorf("zfmtbytes: writing to %T: %w", v, err)
		}
		return nil
	default:
		return fmt.Errorf("zfmtbytes: cannot unmarshal into %T, want *string, *[]byte or io.Writer", v)
	}
}
