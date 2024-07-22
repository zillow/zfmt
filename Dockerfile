FROM golang:1.22 AS build

ENV CGO_ENABLED=0

WORKDIR /go/src/zfmt
COPY . .

RUN go install -v ./...
RUN go build -o zfmt
FROM debian

COPY --from=build /go/src/zfmt /


ENTRYPOINT ["/zfmt"]