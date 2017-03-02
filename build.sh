#!/bin/sh -eux

rm -rf bin/
# rm -rf vendor/

glide install -v

export CGO_ENABLED=0
export GOARCH=amd64

mkdir -p bin
GOOS=darwin  go build -o bin/terraform-provider-kapacitor.macos
GOOS=linux   go build -o bin/terraform-provider-kapacitor.linux
GOOS=windows go build -o bin/terraform-provider-kapacitor.exe
