package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/relaxnow/vcgopkg/files"
	"github.com/relaxnow/vcgopkg/program"
	// "golang.org/x/mod/modfile"
)

func main() {
	flag.Parse()
	inputPath := flag.Arg(0)
	fmt.Println(inputPath)

	inputPathStat, err := os.Stat(inputPath)
	if err != nil {
		log.Printf("Error getting stat for '%s'. File or path may not exist? Original error: '%s'", inputPath, err)
		os.Exit(1)
		return
	}

	log.Print("Reading go files in: " + inputPath)

	featureFiles := files.FeatureFiles{}
	featureFiles.DetectFromPath(inputPathStat)

	if len(featureFiles.MainFiles) == 0 {
		panic("No main files found")
	}

	for _, mainFile := range featureFiles.MainFiles {
		program := program.GetProgramFromMainFilePath(mainFile.FilePath)
		tempProgram := program.CopyToTempDir()
		err := tempProgram.Vendor()
		tempProgram.Zip()
	}

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
