package proto

import "github.com/zillow/zfmt/util"

// SchematizedBase64FormatterDeprecated follows the Confluent Wire Format https://docs.confluent.io/current/schema-registry/serdes-develop/index.html#wire-format
type SchematizedBase64FormatterDeprecated struct {
	formatter Base64Formatter
	SchemaID  int
}

// Marshall converts input into avro binary data with schema ID attached
func (p *SchematizedBase64FormatterDeprecated) Marshall(v any) ([]byte, error) {
	return util.MarshallSchema(&p.formatter, p.SchemaID, v)
}

// Unmarshal fills avro binary data into provided interface v and validates the schema ID
func (p *SchematizedBase64FormatterDeprecated) Unmarshal(b []byte, v any) error {
	return util.UnmarshalSchema(&p.formatter, p.SchemaID, b, v)
}
