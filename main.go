package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/otiai10/copy"
)

func main() {
	flag.Parse()
	inputPath := flag.Arg(0)
	packageDate := flag.Arg(1)
	log.Debug("inputPath=" + inputPath)
	log.Debug("packageDate=" + packageDate)

	absPath, err := filepath.Abs(inputPath)
	if err != nil {
		log.WithFields(log.Fields{
			"inputPath":     inputPath,
			"originalError": err,
		}).Fatal("Error getting absolute path. File or path may not exist")
	}

	absPathStat, err := os.Stat(absPath)
	if err != nil {
		log.WithFields(log.Fields{
			"absPath":       absPath,
			"originalError": err,
		}).Fatal("Error getting stat. File or path may not exist.")
	}

	log.WithField("absPath", "absPath").Debug("Reading go files")

	mainFiles := getMainFiles(absPathStat, absPath)

	log.WithField("mainFiles", mainFiles).Debug("Finished getting mainFiles")

	for _, mainFile := range mainFiles {
		packageMainFile(mainFile, packageDate)
	}
}

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func getMainFiles(absPathStat os.FileInfo, absPath string) []string {
	mainFiles := []string{}
	if absPathStat.IsDir() {
		log.WithField("absPath", absPath).Debug("Input is dir")
		err := filepath.Walk(absPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					return nil
				}

				parsedPackages, err := parser.ParseDir(
					token.NewFileSet(),
					path,
					nil,
					parser.ParseComments,
				)

				if err != nil {
					panic(err)
				}

				for _, pkg := range parsedPackages {
					for _, file := range pkg.Files {
						for _, decl := range file.Decls {
							if ast.FilterDecl(decl, func(name string) bool { return name == "main" }) {
								mainFiles = append(mainFiles, path+"/"+file.Name.Name+".go")
							}
						}
					}
				}
				return nil
			})
		if err != nil {
			log.Fatal(err)
		}
	} else if absPathStat.Mode().Perm().IsRegular() {
		log.Printf("'%s' input is a file", absPath)
		parsedFile, err := parser.ParseFile(token.NewFileSet(), absPath, nil, parser.ParseComments)

		if err != nil {
			log.Fatal(err)
		}

		for _, decl := range parsedFile.Decls {
			if ast.FilterDecl(decl, func(name string) bool { return name == "main" }) {
				mainFiles = append(mainFiles, absPath)
			}
		}
	} else {
		log.Fatalf("'%s' does not exist", absPath)
	}
	if len(mainFiles) == 0 {
		log.Fatalf("No main files found in %s", absPath)
	}
	return mainFiles
}

func packageMainFile(mainFile string, packageDate string) {
	goModPath := ""
	parentDir := filepath.Dir(mainFile)
	log.WithField("parentDir", parentDir).Debug("Starting looking up for mainFile")
	for {
		goModStat, _ := os.Stat(parentDir + "/go.mod")

		if goModStat != nil {
			goModPath = parentDir + "/go.mod"
			log.WithField("goModPath", goModPath).Debug("Found go.mod path")
			break
		}
		if parentDir != "" {
			parentDir = filepath.Dir(parentDir)
			log.WithField("parentDir", parentDir).Debug("Trying parent directory")
			continue
		}

		log.Fatal("go.mod not found")
	}

	tempWorkDir, err := os.MkdirTemp("", "vcgopkg")
	log.WithField("tempWorkDir", tempWorkDir).Debug("Creating temporary working directory")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempWorkDir)

	copyDir := tempWorkDir + "/" + filepath.Base(filepath.Dir(goModPath))

	log.WithFields(log.Fields{"from": parentDir, "to": copyDir}).Debug("Copying files")
	copy.Copy(parentDir, copyDir)

	logFiles(copyDir, "Copied Files")

	vendorDir(copyDir)

	updateVeracodeJson(mainFile, parentDir, copyDir)

	pkg(goModPath, tempWorkDir, parentDir, packageDate)

	logFiles(tempWorkDir, "Temporary workdir after packaging")
	logFiles(parentDir, "ParentDir")
	logFiles(parentDir+"/veracode", "Veracode Dir")
}

func logFiles(dir string, msg string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	fileList := []string{}
	for _, file := range files {
		fileList = append(fileList, file.Name())
	}

	log.WithFields(log.Fields{
		"dir":      dir,
		"fileList": fileList,
	}).Debug(msg)
}

func vendorDir(copyDir string) {
	cmd := exec.Command("go", "mod", "vendor")
	cmd.Dir = copyDir
	cmdOut, _ := cmd.Output()
	log.Debug(string(cmdOut))
}

func updateVeracodeJson(mainFile string, parentDir string, copyDir string) {
	mainFileRelativePath := strings.TrimPrefix(path.Base(mainFile), parentDir)
	json := []byte(fmt.Sprintf("{\"MainFile\": \"%s\"}", mainFileRelativePath))
	ioutil.WriteFile(copyDir+"/veracode.json", json, 0644)
}

func pkg(goModPath string, tempWorkDir string, parentDir string, packageDate string) {
	goModDir := filepath.Dir(goModPath)
	log.Debug(goModDir)
	baseDir := filepath.Base(goModDir)
	if packageDate == "" {
		packageDate = time.Now().Format("-20060102150405")
	}
	zipFile := baseDir + packageDate + ".zip"
	log.Debug(tempWorkDir + "# zip -r " + zipFile + " " + baseDir)
	cmd := exec.Command("zip", "-r", zipFile, baseDir)
	cmd.Dir = tempWorkDir
	cmdOut, _ := cmd.Output()
	log.Debug(string(cmdOut))

	veracodeDir := parentDir + "/veracode"
	os.Mkdir(veracodeDir, 0700)
	log.WithField("veracodeDir", veracodeDir).Debug("Created veracode dir for binaries")

	os.Rename(tempWorkDir+"/"+zipFile, veracodeDir+"/"+zipFile)
	log.WithFields(log.Fields{
		"from": tempWorkDir + "/" + zipFile,
		"to":   veracodeDir + "/" + zipFile,
	}).Debug("Rename zipfile")
}
