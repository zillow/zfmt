# zfmt

[![License](https://img.shields.io/github/license/zillow/zfmt)](https://github.com/zillow/zfmt/blob/main/LICENSE)
[![GitHub Actions](https://github.com/zillow/zfmt/actions/workflows/go.yml/badge.svg)](https://github.com/zillow/zfmt/actions/workflows/go.yml)
[![Codecov](https://codecov.io/gh/zillow/zfmt/branch/main/graph/badge.svg?token=STRT8T67YP)](https://codecov.io/gh/zillow/zfmt)

## About

zfmt is a set of small, independently versioned Go modules, each providing concrete formatter types that
serialize and deserialize Go values. Every formatter has the same shape:

```go
Marshal(v any) ([]byte, error)
Unmarshal(data []byte, v any) error
```

Transport libraries such as [zkafka](https://github.com/zillow/zkafka) and zsqs accept any value with these two
methods. Pick the formatter module you need at the call site and pass an instance in code. You only pull in the
libraries the modules you import need: zfmtjson and zfmtbytes import only the standard library, zfmtproto imports
the protobuf runtime, and zfmtavro imports heetch/avro.

Prior to this split, zfmt was a single module that selected a formatter by name from a string such as `"proto_raw"`.
That registry forced every consumer to depend on every serialization library and turned misconfiguration into a
runtime failure. See [Migrating from 1.x](#migrating-from-1x).

## Modules

| Module                             | Formatters                                                    | Third-party imports          |
|------------------------------------|---------------------------------------------------------------|------------------------------|
| `github.com/zillow/zfmt/zfmtjson`  | `Formatter`                                                   | none                         |
| `github.com/zillow/zfmt/zfmtbytes` | `Formatter`: raw bytes for `string`, `[]byte` and `io.Reader` | none                         |
| `github.com/zillow/zfmt/zfmtproto` | `RawFormatter`, `Base64Formatter`, `JSONFormatter`            | `google.golang.org/protobuf` |
| `github.com/zillow/zfmt/zfmtavro`  | `Formatter`, `SchematizedFormatter`                           | `github.com/heetch/avro`     |

"Third-party imports" counts what the library code imports. Each `go.mod` also lists test-only dependencies such as
testify and go-cmp, so those still appear in a consumer's module graph.

There is no root module. Each formatter module is independent and can be required on its own.

## Installation

Install only the module(s) you need:

```bash
go get github.com/zillow/zfmt/zfmtjson
```

```bash
go get github.com/zillow/zfmt/zfmtproto
```

```bash
go get github.com/zillow/zfmt/zfmtavro
```

```bash
go get github.com/zillow/zfmt/zfmtbytes
```

## Usage

Package names carry a `zfmt` prefix so they never collide with `encoding/json` or
`google.golang.org/protobuf/proto` and need no alias.

```go
import (
	"github.com/zillow/zfmt/zfmtjson"
	"github.com/zillow/zfmt/zfmtproto"
	"github.com/zillow/zfmt/zfmtavro"
)

// encoding/json/v2 with encoding/json (v1) wire behavior
jsonFmt := zfmtjson.Formatter{}

// protobuf binary wire format
protoFmt := zfmtproto.RawFormatter{}

// avro with a Confluent wire format schema ID prefix
avroFmt := zfmtavro.NewSchematizedFormatter(99)

b, err := protoFmt.Marshal(&myMessage)
err = protoFmt.Unmarshal(b, &myMessage)
```

The zfmtjson, zfmtbytes and zfmtproto formatters are ready to use as zero values. The zfmtavro formatters must be
created with `NewFormatter` or `NewSchematizedFormatter`, because each one caches a compiled decoder per record type.
Copies of an avro formatter share that cache.

zfmtavro decodes on the assumption that a message was written with the target type's own schema. It does not
resolve messages written with a different, even compatible, schema. Use a schema registry aware marshaler, such as
zkafka's, when producers and consumers evolve schemas independently.

To implement your own formatter, write a type with the two methods above. No zfmt import is required.

## Migrating from 1.x

In 1.x a formatter was chosen by `zfmt.FormatterType` and `zfmt.GetFormatter`. Both are gone. Replace the string
with the concrete type from the matching module:

| 1.x `FormatterType`       | Replacement                                                           | Module                              |
|---------------------------|-----------------------------------------------------------------------|-------------------------------------|
| `json`                    | `zfmtjson.Formatter{}`                                                | `github.com/zillow/zfmt/zfmtjson`   |
| `string`                  | `zfmtbytes.Formatter{}`                                               | `github.com/zillow/zfmt/zfmtbytes`  |
| `proto_raw`               | `zfmtproto.RawFormatter{}`                                            | `github.com/zillow/zfmt/zfmtproto`  |
| `proto_base64`            | `zfmtproto.Base64Formatter{}`                                         | `github.com/zillow/zfmt/zfmtproto`  |
| `proto_json`              | `zfmtproto.JSONFormatter{}`                                           | `github.com/zillow/zfmt/zfmtproto`  |
| `avro`                    | `zfmtavro.NewFormatter()`                                             | `github.com/zillow/zfmt/zfmtavro`   |
| `avro_schema`             | `zfmtavro.NewSchematizedFormatter(schemaID)`                          | `github.com/zillow/zfmt/zfmtavro`   |
| `json_schema`             | Removed. Use the JSON schema registry marshaler in zkafka.            |                                     |
| `proto_schema_deprecated` | Removed. Use the proto schema registry marshaler in zkafka.           |                                     |

Other breaking changes to plan for:

- `zfmtproto` accepts only messages generated with `google.golang.org/protobuf` (APIv2). Regenerate `.pb.go` files
  built against the deprecated `github.com/golang/protobuf`.
- `zfmtjson` uses `encoding/json/v2` but keeps v1's wire behavior for nil slices and maps and for case-insensitive
  name matching. Invalid UTF-8 and duplicate object names are now errors instead of being silently accepted.
- `zfmtavro` accepts only `github.com/heetch/avro` generated types. Types generated with
  `github.com/actgardner/gogen-avro` are no longer supported; regenerate them with `avrogo`.
- `zfmtbytes` (the 1.x `string` formatter) now also unmarshals into `*string` and `*[]byte`.
- Errors name the module and the failing step and wrap their cause. Code that matched on 1.x error strings needs
  updating.

The 1.x `zfmt.Formatter` interface has no replacement in zfmt. Transport libraries declare the interface they
accept; any formatter here satisfies it.

`github.com/zillow/zfmt` at v1.x remains available and unchanged for consumers that have not migrated.

## Development

Every `make` target loops over all modules and runs with `GOWORK=off`, so each module is built, tested and linted
against its own `go.mod`, as a consumer would resolve it. The `go.work` file is there for gopls and ad hoc commands
across modules.

```bash
make build
```

```bash
make test
```

```bash
make lint
```

```bash
make tidy-check
```

### Releasing

Each module is tagged independently using the Go multi-module convention of a directory prefix:
`zfmtjson/v1.x.y`, `zfmtproto/v1.x.y`, `zfmtavro/v1.x.y` and `zfmtbytes/v1.x.y`. The modules do not depend on
each other, so they can be tagged in any order.
