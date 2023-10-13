#!/bin/bash
set -e

rm -f vcgopkg.exe vcgopkg vcgopkg*.zip

echo "Build for Windows - WinRes"
cd ~
go install github.com/tc-hib/go-winres@latest
cd -
~/go/bin/go-winres make


echo "Build for Windows - Go"
GOOS=windows GOARCH=amd64 go build -o vcgopkg.exe
zip -r vcgopkg-windows-amd64.zip vcgopkg.exe

echo "Build for OS X"
GOOS=darwin GOARCH=amd64 go build -o vcgopkg
zip -r vcgopkg-darwin-amd64.zip vcgopkg
GOOS=darwin GOARCH=arm64 go build -o vcgopkg
zip -r vcgopkg-darwin-arm64.zip vcgopkg

echo "Build for Linux"
GOOS=linux GOARCH=amd64 go build -o vcgopkg
zip -r vcgopkg-linux-amd64.zip vcgopkg
GOOS=linux GOARCH=arm64 go build -o vcgopkg
zip -r vcgopkg-linux-arm64.zip vcgopkg

rm -f vcgopkg.exe vcgopkg