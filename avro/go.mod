module github.com/zillow/zfmt/avro

go 1.22

toolchain go1.22.1

require (
	github.com/actgardner/gogen-avro/v10 v10.2.1
	github.com/heetch/avro v0.4.5
	github.com/zillow/zfmt v0.0.0
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
)

replace github.com/zillow/zfmt => ../.
