FROM golang:1.20.5-bullseye AS builder

WORKDIR /go/src/github.com/hatappi/aws-cloudwatch-golang-filter

ENV GOFLAGS=-buildvcs=false

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o /source/aws-cloudwatch.so -buildmode=c-shared .

FROM alpine

COPY --from=builder /source/aws-cloudwatch.so ./aws-cloudwatch.so
