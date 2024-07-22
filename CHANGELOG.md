# Changelog

All notable changes to this project will be documented in this file.

This project adheres to Semantic Versioning.


## 1.0.0 (September 2022)

    - Enum `ProtoFmt`(`proto`) renamed to `ProtoBase64Fmt`(`proto_base64`). This is intended for SQS use only while the newly introduced `ProtoRawFmt`(`proto_raw`) is intended for most other use cases.
    - Enum `ProtoSchemaFmt`(`proto_schema`) renamed to `ProtoSchemaDeprecatedFmt`(`proto_schema_deprecated`). The proto schema is deprecated because it doesn't work properly with the confluent schema registry. Use the `ProtoRawFmt` or `ProtoBase64Fmt` instead.
