package zfmt

import (
	err "fmt"
)

// FormatterType defines a formatter
type FormatterType string

const (
	// JSONFmt indicates json formatter
	JSONFmt FormatterType = "json"
	// ProtoRawFmt indicates protocol buffer formatter. Does not Base64 encode/decode.
	ProtoRawFmt FormatterType = "proto_raw"
	// ProtoBase64Fmt indicates protocol buffer formatter. Base64 encodes/decodes the message. Only intended for use with SQS.
	ProtoBase64Fmt FormatterType = "proto_base64"
	// ProtoJSONFmt indicates usage protojson formatter (which should be used for json formatted proto messages).
	ProtoJSONFmt FormatterType = "proto_json"
	// StringFmt indicates string formatter
	StringFmt FormatterType = "string"
	// AvroFmt indicates apache avro formatter
	AvroFmt FormatterType = "avro"
	// AvroSchemaFmt indicates apache avro formatter with schemaID encoded
	AvroSchemaFmt FormatterType = "avro_schema"
	// JSONSchemaFmt indicates json formatter with schemaID encoded
	JSONSchemaFmt FormatterType = "json_schema"
	// ProtoSchemaDeprecatedFmt indicates proto formatter with schemaID encoded - deprecated because it doesn't work properly.
	ProtoSchemaDeprecatedFmt FormatterType = "proto_schema_deprecated"
)

// Formatter allows the user to extend formatting capability to unsupported data types
type Formatter interface {
	Marshall(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

// GetFormatter returns supported formatter from its name
func GetFormatter(fmt FormatterType, schemaID int) (Formatter, error) {
	switch fmt {
	case StringFmt:
		return &StringFormatter{}, nil
	case JSONFmt:
		return &JSONFormatter{}, nil
	case ProtoJSONFmt:
		return &ProtoJSONFormatter{}, nil
	case ProtoBase64Fmt:
		return &ProtobufBase64Formatter{}, nil
	case ProtoRawFmt:
		return &ProtobufRawFormatter{}, nil
	case AvroFmt:
		return &AvroFormatter{}, nil
	case AvroSchemaFmt:
		return &SchematizedAvroFormatter{SchemaID: schemaID}, nil
	case JSONSchemaFmt:
		return &SchematizedJSONFormatter{SchemaID: schemaID}, nil
	case ProtoSchemaDeprecatedFmt:
		return &SchematizedProtoFormatterDeprecated{SchemaID: schemaID}, nil
	default:
		return nil, err.Errorf("unsupported formatter %s", fmt)
	}
}
