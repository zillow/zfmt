package zfmt

import (
	"fmt"
	"io"
)

// StringFormatter ...
type StringFormatter struct{}

// Marshall ...
func (f *StringFormatter) Marshall(i any) ([]byte, error) {
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

func (f *StringFormatter) Unmarshal(b []byte, i any) error {
	// since string is immutable, we rely on io.Writer
	switch v := i.(type) {
	case io.Writer:
		_, err := v.Write(b)
		return err
	default:
		return fmt.Errorf("cannot convert from binary %T", i)
	}
}
