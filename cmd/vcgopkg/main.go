package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	// "golang.org/x/mod/modfile"
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
		parsedPackage, err := parser.ParseDir(token.NewFileSet(), dirOrFile, nil, parser.ParseComments)

		if err != nil {
			panic(err)
		}
		spew.Dump(parsedPackage)

	} else if dirOrFileStat.Mode().Perm().IsRegular() {
		log.Printf("'%s' input is a file", dirOrFile)
		parsedFile, err := parser.ParseFile(token.NewFileSet(), dirOrFile, nil, parser.ParseComments)

		if err != nil {
			panic(err)
		}

		for _, decl := range parsedFile.Decls {
			if ast.FilterDecl(decl, func(name string) bool { return name == "main" }) {
				log.Println("Found main!")
				spew.Dump(decl)
			}
		}
	} else {
		log.Printf("'%s' does not exist", dirOrFile)
		os.Exit(1)
	}

	log.Print("Reading go files in: " + dirOrFile)

	// Find all go files, foreach go file get FeatureFiles
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

// func hasGoMod(path) {
// 	mod, err := modfile.Parse(path)
// }

// func isWorkspaceModeWithGoEnv(path) {
// 	// if !GOPATH
// 	//    loadEnvFromGoEnv()
// 	// return isWorkspaceMode(path)
// }

// func isWorkspaceMode(path) {
// 	// workspaceRoot := getWorkspaceRoot()
// 	// if workspaceRoot == nil {
// 	//    return false
// 	// }
// 	// return path.startsWith(path)
// }

// func getWorkspaceRoot() {
// 	// if isWorkspace(GOPATH)
// 	//   return GOPATH
// 	// if isWorkspace(GOROOT)
// 	//    return GOROOT
// 	// return nil
// }

// func isWorkspace(path) {
// 	// if !path.Stat().isDir()
// 	//   return false
// 	// if !path + '/src'.isDir()
// 	//   return false
// 	// return true
// }

// type FeatureFiles struct {
// 	MainFiles      []ast.File
// 	CFiles         []ast.File
// 	mainFiles      []ast.File
// 	importFiles    []ast.File
// 	BuildFiles     []ast.File
// 	OSFiles        []ast.File
// 	FrameworkFiles FrameworkFiles
// }

// type FrameworkFiles struct {
// 	RevelFiles    []ast.File
// 	GinFiles      []ast.File
// 	MartiniFiles  []ast.File
// 	GoWebFiles    []ast.File
// 	GorillaFiles  []ast.File
// 	GojiFiles     []ast.File
// 	GoaFiles      []ast.File
// 	BeegoFiles    []ast.File
// 	BuffaloFiles  []ast.File
// 	kitFiles      []ast.File
// 	echoFiles     []ast.File
// 	fasthttpFiles []ast.File
// 	govwaFiles    []ast.File
// }
