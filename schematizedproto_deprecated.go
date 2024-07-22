package zfmt

// SchematizedProtoFormatterDeprecated follows the Confluent Wire Format https://docs.confluent.io/current/schema-registry/serdes-develop/index.html#wire-format
type SchematizedProtoFormatterDeprecated struct {
	formatter ProtobufBase64Formatter
	SchemaID  int
}

// Marshall converts input into avro binary data with schema ID attached
func (p *SchematizedProtoFormatterDeprecated) Marshall(v interface{}) ([]byte, error) {
	return marshall(&p.formatter, p.SchemaID, v)
}

// Unmarshal fills avro binary data into provided interface v and validates the schema ID
func (p *SchematizedProtoFormatterDeprecated) Unmarshal(b []byte, v interface{}) error {
	return unmarshal(&p.formatter, p.SchemaID, b, v)
}
