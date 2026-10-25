package zfmtproto

import "github.com/zillow/zfmt/v2/schema"

// SchematizedBase64FormatterDeprecated follows the Confluent Wire Format https://docs.confluent.io/current/schema-registry/serdes-develop/index.html#wire-format
//
// Deprecated: this formatter base64 encodes the proto payload before prepending the schema ID, which does not
// interoperate with the confluent schema registry. It is retained only for compatibility with existing data.
// Use RawFormatter, Base64Formatter, or a schema registry aware marshaler in the transport library instead.
type SchematizedBase64FormatterDeprecated struct {
	formatter Base64Formatter
	// SchemaID is prepended to every marshaled payload and validated on unmarshal.
	// A zero SchemaID skips validation on unmarshal.
	SchemaID int
}

// NewSchematizedBase64FormatterDeprecated returns a SchematizedBase64FormatterDeprecated for the given schemaID
//
// Deprecated: see SchematizedBase64FormatterDeprecated.
func NewSchematizedBase64FormatterDeprecated(schemaID int) SchematizedBase64FormatterDeprecated {
	return SchematizedBase64FormatterDeprecated{
		formatter: NewBase64Formatter(),
		SchemaID:  schemaID,
	}
}

// Marshal converts input into base64 proto binary data with schema ID attached
func (p SchematizedBase64FormatterDeprecated) Marshal(v any) ([]byte, error) {
	return schema.Marshal(p.formatter, p.SchemaID, v)
}

// Unmarshal fills base64 proto binary data into provided interface v and validates the schema ID
func (p SchematizedBase64FormatterDeprecated) Unmarshal(b []byte, v any) error {
	return schema.Unmarshal(p.formatter, p.SchemaID, b, v)
}
