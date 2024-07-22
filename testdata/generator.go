package avro

//go:generate $GOPATH/bin/gogen-avro -containers . ./example.avsc
//go:generate protoc --proto_path=. --go_out=./ ./example.proto
//go:generate $GOPATH/bin/avrogo -p heetch -d ./heetch ./example.avsc
