# vcgopkg
Unofficial Community Project to help package a Go application for Veracode Static Analysis

## Linux & OS X

### Download

On Linux with go get:
```
export GOPATH=`go env GOPATH` &&
export PATH="$GOPATH/bin:$PATH" &&
go install github.com/relaxnow/vcgopkg
```

### Usage

Package the current working directory:
```
~/Projects/my-go-project# vcgopkg
```

OR package a directory:
```
~/# vcgopkg Projects/my-go-project
```

OR package a main.go file
```
~/# vcgopkg Projects/my-go-project/cmd/main.go
```
vcgopkg will then look for all main funcs and produce a .zip file for each, for example:

```
~/Projects/my-go-project/veracode/my-go-project--cmd--20210909010101.zip
```

All .zip files from veracode can then be uploaded to Veracode for Static Analysis.

## Windows

### Download

[Download vcgopkg-amd64.exe](https://github.com/relaxnow/vcgopkg/releases/download/v0.0.11/vcgopkg-amd64.exe).

### Usage

Drop the exe inside the go project and double click.

OR package the current working directory with the command line:
```
C:\my-go-project> vcgopkg-amd64
```

OR package a directory:
```
C:\> vcgopkg-amd64 my-go-project
```

OR package a main.go file
```
C:\my-go-project> vcgopkg-amd64 my-go-project\cmd\main.go
```
vcgopkg will then look for all main funcs and produce a .zip file for each, for example:

C:\my-go-project\veracode\my-go-project--cmd--20210909010101.zip

All .zip files from veracode can then be uploaded to Veracode for Static Analysis.
