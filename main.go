package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/otiai10/copy"
)

// TODO: Show version information
// TODO: implement non-debug mode
// TODO: implement versioning
// TODO: implement update check
// TODO: implement help
// TODO: package log inside zip files.
// TODO: vendor only once per go module
// TODO: Detect and show go version
// TODO: Detect and warn on incorrect Go version based on go mod
// TODO: Better error handling when go mod vendor fails
func main() {
	log.Debug("Running version v0.0.12")

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
		}).Panic("Error getting absolute path. File or path may not exist")
		panic("Error getting absolute path. File or path may not exist")
	}
	log.WithField("absPath", absPath).Debug("Reading go files")

	absPathStat, err := os.Stat(absPath)
	if err != nil {
		log.WithFields(log.Fields{
			"absPath":       absPath,
			"originalError": err,
		}).Panic("Error getting stat. File or path may not exist.")
		panic("Error getting stat. File or path may not exist.")
	}

	mainFiles, err := getMainFiles(absPathStat, absPath)
	if err != nil {
		log.WithFields(log.Fields{
			"inputPath":     inputPath,
			"originalError": err,
		}).Panic("Error getting main files")
		panic("Error getting main files")
	}

	log.WithField("mainFiles", mainFiles).Debug("Finished getting mainFiles")

	for _, mainFile := range mainFiles {
		log.WithField("MainFile", mainFile).Debug("Packaging for mainfile")
		err = packageMainFile(mainFile, packageDate)
		if err != nil {
			log.WithFields(log.Fields{
				"mainFile":      mainFile,
				"originalError": err,
			}).Panic("Error getting main file")
			panic("Error getting main file")
		}
	}

	log.Debug("Ran version v0.0.12")
}

func getMainFiles(absPathStat os.FileInfo, absPath string) ([]string, error) {
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
					return err
				}

				for _, pkg := range parsedPackages {
					for filename, file := range pkg.Files {
						for _, decl := range file.Decls {
							if ast.FilterDecl(decl, func(name string) bool { return name == "main" }) {
								mainFiles = append(mainFiles, filename)
							}
						}
					}
				}
				return nil
			})
		if err != nil {
			return []string{}, err
		}
	} else if absPathStat.Mode().Perm().IsRegular() {
		log.Printf("'%s' input is a file", absPath)
		parsedFile, err := parser.ParseFile(token.NewFileSet(), absPath, nil, parser.ParseComments)

		if err != nil {
			return []string{}, err
		}

		for _, decl := range parsedFile.Decls {
			if ast.FilterDecl(decl, func(name string) bool { return name == "main" }) {
				mainFiles = append(mainFiles, absPath)
			}
		}
	} else {
		return []string{}, fmt.Errorf("'%s' does not exist", absPath)
	}
	if len(mainFiles) == 0 {
		return []string{}, fmt.Errorf("no main files found in %s", absPath)
	}
	return mainFiles, nil
}

// TODO: Make work with GOPATH
func packageMainFile(mainFile string, packageDate string) error {
	goModPath := ""
	parentDir := filepath.Dir(mainFile)
	prevParentDir := ""
	log.WithField("parentDir", parentDir).Debug("Starting looking up for go.mod " + mainFile)
	for {
		goModStat, _ := os.Stat(parentDir + "/go.mod")

		if goModStat != nil {
			goModPath = parentDir + "/go.mod"
			log.WithField("goModPath", goModPath).Debug("Found go.mod path")
			break
		}
		if parentDir != prevParentDir {
			prevParentDir = parentDir
			parentDir = filepath.Dir(parentDir)
			log.WithField("parentDir", parentDir).Debug("Trying parent directory")
			continue
		}

		return errors.New("go.mod not found")
	}

	tempWorkDir, err := os.MkdirTemp("", "vcgopkg")
	log.WithField("tempWorkDir", tempWorkDir).Debug("Creating temporary working directory")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempWorkDir)

	copyDir := tempWorkDir + "/" + filepath.Base(filepath.Dir(goModPath))

	log.WithFields(log.Fields{"from": parentDir, "to": copyDir}).Debug("Copying files")
	err = copy.Copy(parentDir, copyDir)
	if err != nil {
		return err
	}

	// TODO: Don't copy veracode dir, instead of copying it and then removing it
	log.WithField("dir", copyDir+"/veracode").Debug("Removing veracode directory from copy")
	err = os.RemoveAll(copyDir + "/veracode")
	if err != nil {
		return err
	}

	LogFiles(copyDir, "Copied Files")

	err = vendorDir(copyDir)
	if err != nil {
		return err
	}

	err = updateVeracodeJson(mainFile, parentDir, copyDir)
	if err != nil {
		return err
	}

	err = pkg(goModPath, mainFile, tempWorkDir, parentDir, packageDate)
	if err != nil {
		return err
	}

	LogFiles(tempWorkDir, "Temporary workdir after packaging")
	LogFiles(parentDir, "ParentDir")
	LogFiles(parentDir+"/veracode", "Veracode Dir")
	return nil
}

func vendorDir(copyDir string) error {
	log.Debug("Starting vendor")

	_, goPathErr := exec.LookPath("go")
	log.Debug(goPathErr)
	_, vendorPathErr := os.Stat(copyDir + "/vendor")
	log.WithField("stat error", vendorPathErr).Debug("Tested if vendor directory exists")

	canFindVendor := vendorPathErr == nil
	canFindGoExecutable := goPathErr == nil

	if !canFindVendor && !canFindGoExecutable {
		log.Debug("Error getting stat of /vendor")
		return fmt.Errorf("no vendor directory and unable to run go mod vendor")
	} else if !canFindGoExecutable && canFindVendor {
		log.Debug("No go but an existing vendor dir, chancing it")
		return nil
	}

	log.Debug("Vendor folder did not exist, running go mod vendor, this may take a while")
	cmd := exec.Command("go", "mod", "vendor")
	cmd.Dir = copyDir
	cmdOut, err := cmd.Output()
	log.WithFields(log.Fields{
		"cmdOut": string(cmdOut),
	}).Debug("Ran go mod vendor")
	return err
}

// TODO: Find FirstParty
func updateVeracodeJson(mainFile string, parentDir string, copyDir string) error {
	log.Debug("Updating veracode.json")
	file := copyDir + "/" + "veracode.json"
	err := CreateEmptyVeracodeJsonFileIfNotExists(file)

	if err != nil {
		return err
	}

	veracodeJsonFile, err := NewVeracodeJsonFile(file)

	if err != nil {
		return err
	}

	mainFileRelativePath := strings.TrimPrefix(filepath.Dir(mainFile), parentDir+"/")
	veracodeJsonFile.VeracodeJson.MainRoot = mainFileRelativePath

	return veracodeJsonFile.WriteToFile()
}

// TODO: Allow writing to output directory
// TODO: Use path to main in output file to support multiple path
// TODO: Test package with go loader
func pkg(goModPath string, mainFile string, tempWorkDir string, parentDir string, packageDate string) error {
	goModDir := filepath.Dir(goModPath)
	log.Debug(goModDir)
	log.Debug(mainFile)
	baseDir := filepath.Base(goModDir)
	if packageDate == "" {
		packageDate = time.Now().Format("_20060102150405")
	}
	// Turn /path/to/module/cmd/main.go into cmd-main
	// TODO: use allow-list instead of deny-list
	relativeMainPath := strings.TrimSuffix(mainFile[len(goModDir)+1:], ".go")
	cmdSlug := "_"
	cmdSlug += strings.ReplaceAll(relativeMainPath, "\\", "--")
	cmdSlug = strings.ReplaceAll(cmdSlug, "/", "--")
	zipFile := baseDir + cmdSlug + packageDate + ".zip"
	log.WithFields(log.Fields{
		"baseDir": baseDir,
		"zipFile": tempWorkDir + "/" + zipFile,
	}).Debug("Writing zip file")
	err := ZipWriter(tempWorkDir, tempWorkDir+"/"+zipFile)
	if err != nil {
		return err
	}

	veracodeDir := parentDir + "/" + "veracode"
	_, err = os.Stat(veracodeDir)
	if err != nil {
		err = os.Mkdir(veracodeDir, 0700)
		log.WithField("veracodeDir", veracodeDir).Debug("Created veracode dir for binaries")
		if err != nil {
			return err
		}
	}

	err = MoveFile(tempWorkDir+"/"+zipFile, veracodeDir+"/"+zipFile)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"from": tempWorkDir + "/" + zipFile,
		"to":   veracodeDir + "/" + zipFile,
	}).Debug("Rename zipfile")
	return nil
}
