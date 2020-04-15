#!/bin/bash
echo Building...
GOARCH=amd64 go build -ldflags="-s -w" -o ./build/pwndrop -mod=vendor main.go
echo Done.
