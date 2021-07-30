# vcgopkg
Unofficial Community Project to help package a Go application for Veracode Static Analysis

# IN DEVELOPMENT, DO NOT USE

## Usage

Package the current working directory:
```
repo# vcgopkg
```

Package a directory:
```
vcgopkg repo/cmd
```

Package a main.go file
```
vcgopkg repo/path/to/dir/cmd/main.go
```
vcgopkg will then look for all main funcs and produce a .zip file for each, for example: **repo.zip**. 
This can then be uploaded to Veracode for Static Analysis.

## Download

On Linux with go get:
```
export GOPATH=`go env GOPATH` &&
export PATH="$PATH:$GOPATH/bin" &&
go install github.com/relaxnow/vcgopkg
```