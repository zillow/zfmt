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
methods. Pick the formatter module you need at the call site and pass an instance in code. You only pay for the
dependencies of the modules you import: the zfmtjson and zfmtstring modules use only the standard library, the zfmtproto module
pulls in the protobuf runtime, and the zfmtavro module pulls in the avro libraries.

Prior to 2.0, zfmt was a single module that selected a formatter by name from a string such as `"proto_raw"`.
That registry forced every consumer to depend on every serialization library and turned misconfiguration into a
runtime failure. See [Migrating from 1.x](#migrating-from-1x).

## Modules

| Module                         | Formatters                                                                 | Notable dependencies                                  |
|--------------------------------|----------------------------------------------------------------------------|-------------------------------------------------------|
| `github.com/zillow/zfmt/zfmtjson`  | `Formatter`, `SchematizedFormatter`                                        | standard library only                                 |
| `github.com/zillow/zfmt/zfmtstring` | `StringFormatter`: pass-through for `string`, `[]byte` and `io.Reader` values | standard library only                               |
| `github.com/zillow/zfmt/zfmtproto` | `RawFormatter`, `Base64Formatter`, `JSONFormatter`, `SchematizedBase64FormatterDeprecated` | `google.golang.org/protobuf`, `github.com/golang/protobuf` |
| `github.com/zillow/zfmt/zfmtavro`  | `Formatter`, `SchematizedFormatter`                                        | `github.com/actgardner/gogen-avro`, `github.com/heetch/avro` |
| `github.com/zillow/zfmt/v2`    | `schema` package: Confluent wire format helpers shared by the schematized formatters | standard library only                        |

Service code normally imports one or more of the formatter modules. The root module is an implementation detail of the
schematized formatters, though it can be imported directly to build a custom schematized formatter.

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
go get github.com/zillow/zfmt/zfmtstring
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

// std-lib json
jsonFmt := zfmtjson.NewFormatter()

// protobuf wire format for v1 or v2 generated messages
protoFmt := zfmtproto.NewRawFormatter()

// avro with a Confluent wire format schema ID prefix
avroFmt := zfmtavro.NewSchematizedFormatter(99)

b, err := protoFmt.Marshal(&myMessage)
err = protoFmt.Unmarshal(b, &myMessage)
```

All formatters are stateless value types with value receivers. `zfmtjson.Formatter{}` and `&zfmtjson.Formatter{}`
are both valid wherever a formatter is accepted.

To implement your own formatter, write a type with the two methods above. No zfmt import is required.

## Migrating from 1.x

In 1.x a formatter was chosen by `zfmt.FormatterType` and `zfmt.GetFormatter`. Both are gone. Replace the string
with the concrete type from the matching module:

| 1.x `FormatterType`         | 2.x type                                               | Module                         |
|-----------------------------|--------------------------------------------------------|--------------------------------|
| `json`                      | `zfmtjson.Formatter`                                       | `github.com/zillow/zfmt/zfmtjson`  |
| `string`                    | `zfmtstring.StringFormatter`                               | `github.com/zillow/zfmt/zfmtstring` |
| `json_schema`               | `zfmtjson.SchematizedFormatter` (`NewSchematizedFormatter(schemaID)`) | `github.com/zillow/zfmt/zfmtjson` |
| `proto_raw`                 | `zfmtproto.RawFormatter`                                   | `github.com/zillow/zfmt/zfmtproto` |
| `proto_base64`              | `zfmtproto.Base64Formatter`                                | `github.com/zillow/zfmt/zfmtproto` |
| `proto_json`                | `zfmtproto.JSONFormatter`                                  | `github.com/zillow/zfmt/zfmtproto` |
| `proto_schema_deprecated`   | `zfmtproto.SchematizedBase64FormatterDeprecated` (`NewSchematizedBase64FormatterDeprecated(schemaID)`) | `github.com/zillow/zfmt/zfmtproto` |
| `avro`                      | `zfmtavro.Formatter`                                       | `github.com/zillow/zfmt/zfmtavro`  |
| `avro_schema`               | `zfmtavro.SchematizedFormatter` (`NewSchematizedFormatter(schemaID)`) | `github.com/zillow/zfmt/zfmtavro` |

The 1.x `zfmt.Formatter` interface has no replacement in zfmt. Transport libraries declare the interface they
accept; any 2.x formatter satisfies it.

`github.com/zillow/zfmt` at v1.x remains available and unchanged for consumers that have not migrated.

## Development

The repository is a Go workspace (`go.work`) so the formatter modules build against the root module in the same
checkout. Every `make` target loops over all modules:

```bash
make build
```

```bash
make test
```

```bash
make lint
```

### Releasing

Each module is tagged independently using the Go multi-module convention of a directory prefix:
`v2.x.y` for the root module and `zfmtjson/v1.x.y`, `zfmtproto/v1.x.y`, `zfmtavro/v1.x.y`, `zfmtstring/v1.x.y`
for the formatter modules.

When the root module changes, tag it first and bump the `github.com/zillow/zfmt/v2` requirement in each formatter
module before tagging those. Until the root tag exists, `go mod tidy` inside a formatter module cannot resolve the
requirement; `go build` and `go test` continue to work through the workspace.
