package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	flag.Parse()
	dirOrFile := flag.Arg(0)
	fmt.Println(dirOrFile)

	log.Print("Reading go files in: " + dirOrFile)
	// Find all go files, foreach go file get:
	//	  main
	//    import "C"
	//    Build tags
	//    OS specific features
	//    Framework import
	//       Revel
	//		 Gin
	//       Martini
	//       Web.Go
	//       Gorilla
	//       Goji
	//       Goa
	//       Beego
	//       Buffalo
	//       kit
	//       echo
	//       kit
	//       fasthttp
	//       govwa
	// If no main funcs, error out

	// Foreach main funcs
	//    Get repo root
	//        copy
	//        vendor deps
	//

	// Testcases:
	//   Go multi repo: https://github.com/flowerinthenight/golang-monorepo
	//   GOROOT: https://golang.org/doc/gopath_code
	//   Broken code
	//   Missing imports

}
