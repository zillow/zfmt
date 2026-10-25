package zfmtstring

import (
	"fmt"
	"io"
)

// StringFormatter is a pass-through formatter. It passes strings, byte slices and io.Readers through as raw bytes; Unmarshal writes into an io.Writer
type StringFormatter struct{}

// NewStringFormatter returns a StringFormatter
func NewStringFormatter() StringFormatter {
	return StringFormatter{}
}

// Marshal returns the bytes of a string, []byte or io.Reader
func (f StringFormatter) Marshal(i any) ([]byte, error) {
	switch v := i.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	case io.Reader:
		return io.ReadAll(v)
	default:
		return nil, fmt.Errorf("cannot convert to string %T", i)
	}
}

// Unmarshal writes b into i, which must be an io.Writer
func (f StringFormatter) Unmarshal(b []byte, i any) error {
	// since string is immutable, we rely on io.Writer
	switch v := i.(type) {
	case io.Writer:
		_, err := v.Write(b)
		return err
	default:
		return fmt.Errorf("cannot convert from binary %T", i)
	}
}
