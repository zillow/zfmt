package zfmtjson

import "github.com/zillow/zfmt/v2/schema"

// SchematizedFormatter follows the Confluent Wire Format https://docs.confluent.io/current/schema-registry/serdes-develop/index.html#wire-format
type SchematizedFormatter struct {
	formatter Formatter
	// SchemaID is prepended to every marshaled payload and validated on unmarshal.
	// A zero SchemaID skips validation on unmarshal.
	SchemaID int
}

// NewSchematizedFormatter returns a SchematizedFormatter which encodes schemaID using the Confluent wire format
func NewSchematizedFormatter(schemaID int) SchematizedFormatter {
	return SchematizedFormatter{
		formatter: NewFormatter(),
		SchemaID:  schemaID,
	}
}

// Marshal converts input into json data with schema ID attached
func (p SchematizedFormatter) Marshal(v any) ([]byte, error) {
	return schema.Marshal(p.formatter, p.SchemaID, v)
}

// Unmarshal fills json data into provided interface v and validates the schema ID
func (p SchematizedFormatter) Unmarshal(b []byte, v any) error {
	return schema.Unmarshal(p.formatter, p.SchemaID, b, v)
}
