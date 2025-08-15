package avro

import "github.com/zillow/zfmt/util"

// SchematizedFormatter follows the Confluent Wire Format https://docs.confluent.io/current/schema-registry/serdes-develop/index.html#wire-format
type SchematizedFormatter struct {
	formatter Formatter
	SchemaID  int
}

// Marshall converts input into avro binary data with schema ID attached
func (p *SchematizedFormatter) Marshall(v any) ([]byte, error) {
	return util.MarshallSchema(&p.formatter, p.SchemaID, v)
}

// Unmarshal fills avro binary data into provided interface v and validates the schema ID
func (p *SchematizedFormatter) Unmarshal(b []byte, v any) error {
	return util.UnmarshalSchema(&p.formatter, p.SchemaID, b, v)
}
