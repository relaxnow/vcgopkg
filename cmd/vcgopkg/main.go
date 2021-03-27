package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Parse()
	dirOrFile := flag.Arg(0)
	fmt.Println(dirOrFile)

	dirOrFileStat, err := os.Stat(dirOrFile)
	if err != nil {
		log.Printf("Error getting stat for '%s'. File or path may not exist? Original error: '%s'", dirOrFile, err)
		os.Exit(1)
		return
	}

	if dirOrFileStat.IsDir() {
		log.Printf("'%s' input is dir", dirOrFile)
	} else if dirOrFileStat.Mode().Perm().IsRegular() {
		log.Printf("'%s' input is a file", dirOrFile)
	} else {
		log.Printf("'%s' does not exist", dirOrFile)
	}

	log.Print("Reading go files in: " + dirOrFile)

	// Find all go files, foreach go file get:
	//	  main
	//    import "C"  jyjy
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
	//   Bazel
	//   Broken code
	//   Missing imports
	//	 Windows machine without go installed

	// if root
	//  copy program to temp dir
	//  vendor / go mod vendor
	// else if isWorkspaceModeWithGoEnv(path) {
	//	get all .go files from program
	//  Find all imports
	//  find program root by looking for shortest import that overlaps with cwd
	//  copy program to temp dir
	//  vendor all imports
	// } else if mod := hasGoMod(path) {
	//  copy module to temp dir
	//	go mod vendor
	// } else {
	//	 error
	// }
	//
	//
	//
	//
	//
	//
}

func isWorkspaceModeWithGoEnv(path) {
	// if !GOPATH
	//    loadEnvFromGoEnv()
	// return isWorkspaceMode(path)
}

func isWorkspaceMode(path) {
	// workspaceRoot := getWorkspaceRoot()
	// if workspaceRoot == nil {
	//    return false
	// }
	// return path.startsWith(path)
}

func getWorkspaceRoot() {
	// if isWorkspace(GOPATH)
	//   return GOPATH
	// if isWorkspace(GOROOT)
	//    return GOROOT
	// return nil
}

func isWorkspace(path) {
	// if !path.Stat().isDir()
	//   return false
	// if !path + '/src'.isDir()
	//   return false
	// return true
}
