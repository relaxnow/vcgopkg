package files

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
)

type ParsedFile struct {
	File     *ast.File
	FilePath string
}

type ParsedPackage struct {
	Package  *ast.Package
	FilePath string
}

type FeatureFiles struct {
	MainFiles []ParsedFile
	// CFiles         []ast.File
	// importFiles    []ast.File
	// BuildFiles     []ast.File
	// OSFiles        []ast.File
	// FrameworkFiles FrameworkFiles
}

func (f *FeatureFiles) DetectFromPath(inputPath os.FileInfo) {
	if inputPathStat.IsDir() {
		log.Printf("'%s' input is dir", inputPath)
		err := filepath.Walk(inputPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					return nil
				}

				parsedPackage, err := parser.ParseDir(
					token.NewFileSet(),
					path,
					nil,
					parser.ParseComments,
				)

				if err != nil {
					panic(err)
				}

				for _, pkg := range parsedPackage {
					f.detectFromPackage(ParsedPackage{
						Package:  pkg,
						FilePath: path,
					})
				}
				return nil
			})
		if err != nil {
			log.Fatal(err)
		}

		if len(f.MainFiles) == 0 {
			log.Fatalf("No main files found in %s", inputPath)
		}

		spew.Dump(f)

	} else if inputPathStat.Mode().Perm().IsRegular() {
		log.Printf("'%s' input is a file", inputPath)
		parsedFile, err := parser.ParseFile(token.NewFileSet(), inputPath, nil, parser.ParseComments)

		if err != nil {
			panic(err)
		}

		f.detectFromFile(ParsedFile{File: parsedFile, FilePath: inputPath})
		spew.Dump(f)
	} else {
		log.Fatalf("'%s' does not exist", inputPath)
		os.Exit(1)
	}
}

func (f *FeatureFiles) detectFromFile(parsedFile ParsedFile) {
	for _, decl := range parsedFile.File.Decls {
		if ast.FilterDecl(decl, func(name string) bool { return name == "main" }) {
			f.MainFiles = append(f.MainFiles, parsedFile)
		}
	}
}

func (f *FeatureFiles) detectFromPackage(pkg ParsedPackage) {
	for _, file := range pkg.Package.Files {
		f.detectFromFile(ParsedFile{
			File:     file,
			FilePath: pkg.FilePath + "/" + file.Name.Name + ".go",
		})
	}
}

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
