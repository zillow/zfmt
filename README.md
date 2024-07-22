# zfmt

[![coverage report](https://gitlab.zgtools.net/devex/archetypes/gomods/zfmt/badges/main/coverage.svg)](https://gitlab.zgtools.net/devex/archetypes/gomods/zfmt/-/commits/main)

[![pipeline status](https://gitlab.zgtools.net/devex/archetypes/gomods/zfmt/badges/main/pipeline.svg)](https://gitlab.zgtools.net/devex/archetypes/gomods/zfmt/-/commits/main)


A go library which contains a number of useful implementations of the Formatter interface,
an interface which contains Marshall and Unmarshall methods.

## Dependencies

### Gogen Avro

Install [Gogen Avro](https://github.com/actgardner/gogen-avro) a utility for generating go code from avro schemas (used for testdata)

```
go get github.com/actgardner/gogen-avro/v7/cmd/...
```


### Migration Guide

1. #### V0 to V1 Migration

    * Enum `ProtoFmt`(`proto`) renamed to `ProtoBase64Fmt`(`proto_base64`). This is intended for SQS use only while the newly introduced `ProtoRawFmt`(`proto_raw`) is intended for most other use cases.
    * Enum `ProtoSchemaFmt`(`proto_schema`) renamed to `ProtoSchemaDeprecatedFmt`(`proto_schema_deprecated`). The proto schema is deprecated because it doesn't work properly with the confluent schema registry. Use the `ProtoRawFmt` or `ProtoBase64Fmt` instead.