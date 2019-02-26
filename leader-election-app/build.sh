#!/bin/bash
set -ex

mkdir -p bin
export GOOS=linux
export GOARCH=amd64
go get ./...
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bin/app app.go

# Build image
docker build -t "mboye/leader-election-app:v1" .
