package zfmt

import (
	"fmt"
	"io"
	"io/ioutil"
)

// StringFormatter ...
type StringFormatter struct{}

// Marshall ...
func (f *StringFormatter) Marshall(i interface{}) ([]byte, error) {
	switch v := i.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	case io.Reader:
		return ioutil.ReadAll(v)
	default:
		return nil, fmt.Errorf("cannot convert to string %T", i)
	}
}

// Unmarshal ...
func (f *StringFormatter) Unmarshal(b []byte, i interface{}) error {
	// since string is immutable, we rely on io.Writer
	switch v := i.(type) {
	case io.Writer:
		_, err := v.Write(b)
		return err
	default:
		return fmt.Errorf("cannot convert from binary %T", i)
	}
}
