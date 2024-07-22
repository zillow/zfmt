package zfmt

// SchematizedJSONFormatter follows the Confluent Wire Format https://docs.confluent.io/current/schema-registry/serdes-develop/index.html#wire-format
type SchematizedJSONFormatter struct {
	formatter JSONFormatter
	SchemaID  int
}

// Marshall converts input into avro binary data with schema ID attached
func (p *SchematizedJSONFormatter) Marshall(v any) ([]byte, error) {
	return marshall(&p.formatter, p.SchemaID, v)
}

// Unmarshal fills avro binary data into provided interface v and validates the schema ID
func (p *SchematizedJSONFormatter) Unmarshal(b []byte, v any) error {
	return unmarshal(&p.formatter, p.SchemaID, b, v)
}
