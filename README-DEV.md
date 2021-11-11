# vcgopkg

Guiding philosophy of vcgopkg is to:
* Assume the best; warn instead of error if possible
* Be humble; verbose on error, quiet on success
* Debug log all the things, to better support customers running the application

# TODO

### Options

-o --output <path/to/output.zip>

## Features

### vendoring
### Unsupported feature detection
* CGO
* Build tags
* Plugins
* OS specific features
* Frameworks
* Multi-module

### Multi-main
### veracode.json generation

* Use CFG to restrict scope of warnings
* vcgopkg.log.json ?
* Git commit?

# Testcases:

* Go multi repo: https://github.com/flowerinthenight/golang-monorepo
* GOROOT: https://golang.org/doc/gopath_code
* Bazel
* Broken code
* Missing imports
* Windows machine without go installed
