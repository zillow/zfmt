package zfmtavro

// SchematizedFormatter encodes/decodes heetch/avro generated types in the Confluent wire format
// https://docs.confluent.io/platform/current/schema-registry/fundamentals/serdes-develop/index.html#wire-format
//
// Create it with NewSchematizedFormatter. Like Formatter, it caches a compiled decoder per record type and
// assumes the message was written with the target type's own schema; the schema ID in the header is
// validated, not used to look up a writer schema.
type SchematizedFormatter struct {
	schemaID int
	decoders *decoders
}

// NewSchematizedFormatter returns a SchematizedFormatter that prefixes every marshaled payload with schemaID
// and requires it on unmarshal. A zero schemaID skips the check on unmarshal.
func NewSchematizedFormatter(schemaID int) SchematizedFormatter {
	return SchematizedFormatter{schemaID: schemaID, decoders: newDecoders(true)}
}

// Marshal converts input into avro binary data with schema ID attached
func (s SchematizedFormatter) Marshal(v any) ([]byte, error) {
	payload, err := Formatter{}.Marshal(v)
	if err != nil {
		return nil, err
	}
	return encodeWireFormat(s.schemaID, payload), nil
}

// Unmarshal validates the wire format header and schema ID, then decodes the avro payload into v
func (s SchematizedFormatter) Unmarshal(b []byte, v any) error {
	if s.decoders == nil {
		return errNotConstructed
	}
	if err := checkWireFormat(s.schemaID, b); err != nil {
		return err
	}
	return s.decoders.unmarshal(b, v)
}
