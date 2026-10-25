# Changelog

All notable changes to this project will be documented in this file.

This project adheres to Semantic Versioning. Each module in the repository is versioned independently.

## Unreleased

### Module split

Breaking. The monolithic `github.com/zillow/zfmt` module is split into independently versioned formatter modules so
that consumers only pull in the serialization libraries they use. `github.com/zillow/zfmt` itself stays at v1.0.1
and gets no further releases; there is no root module after the split.

- Removed `FormatterType`, `GetFormatter` and the string-to-formatter registry. Formatters are now constructed in
  code and passed to transport libraries directly.
- Removed the `Formatter` interface. Transport libraries declare the interface they accept.
- Every formatter type moves to one of the modules below. See the README for the name mapping.
- All modules require Go 1.27.
- Dependencies are on their latest releases: `google.golang.org/protobuf` v1.36.12 and `github.com/heetch/avro`
  v0.4.79. No dependency has a newer major version.

### github.com/zillow/zfmt/zfmtjson 1.0.0

- `Formatter` (was `zfmt.JSONFormatter`).
- Encodes with `encoding/json/v2`, configured to keep the wire behavior of `encoding/json` (v1): nil slices and
  maps marshal as `null`, and object names match struct fields case-insensitively without ignoring `_` or `-`.
- Breaking: v2's stricter input checks stay on. Invalid UTF-8 is an error instead of being replaced with U+FFFD, and
  duplicate object names (including names that differ only in case) are an error instead of the last one winning.
- Breaking: `zfmt.SchematizedJSONFormatter` (`json_schema`) is removed. Use the JSON schema registry marshaler in
  zkafka instead.
- Imports only the standard library. Test-only dependencies are listed in `go.mod`.

### github.com/zillow/zfmt/zfmtbytes 1.0.0

- `Formatter` (was `zfmt.StringFormatter`), which passes payloads through as raw bytes.
- `Unmarshal` accepts `*string` and `*[]byte` targets as well as `io.Writer`. 1.x rejected `*string`, the most common
  target. A `*[]byte` target receives a copy of the payload. A nil `*string` or `*[]byte` returns an error.
- Imports only the standard library. Test-only dependencies are listed in `go.mod`.

### github.com/zillow/zfmt/zfmtproto 1.0.0

- `RawFormatter` (was `zfmt.ProtobufRawFormatter`), `Base64Formatter` (was `zfmt.ProtobufBase64Formatter`),
  and `JSONFormatter` (was `zfmt.ProtoJSONFormatter`).
- Breaking: `zfmt.SchematizedProtoFormatterDeprecated` (`proto_schema_deprecated`) is removed. It never
  interoperated with the Confluent schema registry. Use the proto schema registry marshaler in zkafka instead.
- Breaking: drops support for messages generated against the deprecated `github.com/golang/protobuf` (APIv1).
  Only `google.golang.org/protobuf` (APIv2) messages are accepted, and the v1 runtime is no longer a dependency.
  Regenerate older `.pb.go` files with a current `protoc-gen-go`.
- `Unmarshal` returns an error for a typed nil target such as `(*pb.Msg)(nil)`. 1.x panicked inside protobuf.

### github.com/zillow/zfmt/zfmtavro 1.0.0

- `Formatter` (was `zfmt.AvroFormatter`) and `SchematizedFormatter` (was `zfmt.SchematizedAvroFormatter`).
- Breaking: drops support for `github.com/actgardner/gogen-avro` generated types. Only `github.com/heetch/avro`
  generated types are accepted. Regenerate older types with `avrogo`. `gogen-avro` remains a transitive
  dependency because `heetch/avro` uses its schema compiler internally.
- `Formatter.Marshal` accepts a pointer to a generated record and encodes the record it points to, producing the same
  bytes as the value. 1.x panicked inside `heetch/avro` on any such pointer. A nil pointer now returns an error.
- Breaking: create formatters with `NewFormatter` and `NewSchematizedFormatter`. Each formatter caches a compiled
  decoder per record type (copies share the cache), which makes `Unmarshal` about 24x faster with 3 allocations
  instead of 122. A zero-value formatter returns an error from `Unmarshal`. `SchematizedFormatter.SchemaID` is no
  longer an exported field; pass the ID to the constructor.
- `Unmarshal` assumes the message was written with the target type's own schema, as 1.x did. Messages written with
  a different, even compatible, schema are not resolved; use a schema registry aware marshaler for that.
- `Unmarshal` returns an error for a non-pointer or nil target instead of panicking.
- Breaking: `SchematizedFormatter.Unmarshal` rejects payloads whose first byte is not the Confluent magic byte `0`.
  1.x skipped that check, so a payload with no header could decode with wrong values when `SchemaID` was `0`.
- The Confluent wire format helpers are now unexported implementation details of this module.

### All formatter modules

- The zfmtjson, zfmtbytes and zfmtproto formatters are ready to use as zero values and have no constructors. The
  zfmtavro formatters are created with `NewFormatter` and `NewSchematizedFormatter`, because they cache decoders.
- All formatters use value receivers, so a formatter value and a pointer to it are interchangeable.
- Errors carry the module name and the failing step (for example `zfmtproto: decoding base64: ...`) and wrap the
  underlying error with `%w`, so `errors.Is` and `errors.As` see the cause. Code that matched on 1.x error strings
  needs updating.

## 1.0.1 (July 2024)

Updated go.mod path (`gitlab.zgtools.net/devex/archetypes/gomods/zfmt` -> `github.com/zillow/zfmt`)

## 1.0.0 (July 2024)

Initial release to public github.com

Internal Releases Below
----

## 1.0.0 (September 2022)

    - Enum `ProtoFmt`(`proto`) renamed to `ProtoBase64Fmt`(`proto_base64`). This is intended for SQS use only while the newly introduced `ProtoRawFmt`(`proto_raw`) is intended for most other use cases.
    - Enum `ProtoSchemaFmt`(`proto_schema`) renamed to `ProtoSchemaDeprecatedFmt`(`proto_schema_deprecated`). The proto schema is deprecated because it doesn't work properly with the confluent schema registry. Use the `ProtoRawFmt` or `ProtoBase64Fmt` instead.
