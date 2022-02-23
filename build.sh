#!/bin/bash
go install github.com/tc-hib/go-winres@latest
go-winres make
GOOS=windows GOARCH=amd64 go build -o vcgopkg-amd64.exe
