# Changelog

All notable changes to this project will be documented in this file.

This project adheres to Semantic Versioning. Each module in the repository is versioned independently.

## Unreleased

### github.com/zillow/zfmt/v2 2.0.0

Breaking. The monolithic module is split into independently versioned modules so that consumers only pull in the
serialization libraries they use.

- Removed `FormatterType`, `GetFormatter` and the string-to-formatter registry. Formatters are now constructed in
  code and passed to transport libraries directly.
- Removed the `Formatter` interface. Transport libraries declare the interface they accept.
- Removed every formatter type from the root module. They live in the new modules below. See the README for the
  name mapping.
- Added package `schema` with `Marshal` and `Unmarshal` helpers implementing the Confluent wire format
  (0x00 magic byte plus a 4 byte big-endian schema ID), shared by the schematized formatters.
- The root module no longer depends on any third-party library.
- All modules require Go 1.26.

### github.com/zillow/zfmt/zfmtjson 1.0.0

- `Formatter` (was `zfmt.JSONFormatter`) and `SchematizedFormatter` (was `zfmt.SchematizedJSONFormatter`).
- Standard library only.

### github.com/zillow/zfmt/zfmtstring 1.0.0

- `StringFormatter` (was `zfmt.StringFormatter`), a pass-through for `string`, `[]byte` and `io.Reader` values.
- Standard library only, with no dependency on the root module.

### github.com/zillow/zfmt/zfmtproto 1.0.0

- `RawFormatter` (was `zfmt.ProtobufRawFormatter`), `Base64Formatter` (was `zfmt.ProtobufBase64Formatter`),
  `JSONFormatter` (was `zfmt.ProtoJSONFormatter`) and `SchematizedBase64FormatterDeprecated`
  (was `zfmt.SchematizedProtoFormatterDeprecated`).
- Still supports both `github.com/golang/protobuf` v1 and `google.golang.org/protobuf` v2 generated messages.
  This is the only module that depends on the deprecated v1 runtime.

### github.com/zillow/zfmt/zfmtavro 1.0.0

- `Formatter` (was `zfmt.AvroFormatter`) and `SchematizedFormatter` (was `zfmt.SchematizedAvroFormatter`).

### All formatter modules

- Every formatter has a `New...` constructor and uses value receivers, so values and pointers are interchangeable.
- Formatter implementations are otherwise unchanged from 1.0.1.

## 1.0.1 (July 2024)

Updated go.mod path (`gitlab.zgtools.net/devex/archetypes/gomods/zfmt` -> `github.com/zillow/zfmt`)

## 1.0.0 (July 2024)

Initial release to public github.com

Internal Releases Below
----

## 1.0.0 (September 2022)

    - Enum `ProtoFmt`(`proto`) renamed to `ProtoBase64Fmt`(`proto_base64`). This is intended for SQS use only while the newly introduced `ProtoRawFmt`(`proto_raw`) is intended for most other use cases.
    - Enum `ProtoSchemaFmt`(`proto_schema`) renamed to `ProtoSchemaDeprecatedFmt`(`proto_schema_deprecated`). The proto schema is deprecated because it doesn't work properly with the confluent schema registry. Use the `ProtoRawFmt` or `ProtoBase64Fmt` instead.
