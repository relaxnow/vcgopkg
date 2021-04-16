package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
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
								mainFiles = append(mainFiles, path+string(filepath.Separator)+file.Name.Name+".go")
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
		goModStat, _ := os.Stat(parentDir + string(filepath.Separator) + "go.mod")

		if goModStat != nil {
			goModPath = parentDir + string(filepath.Separator) + "go.mod"
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

	copyDir := tempWorkDir + string(filepath.Separator) + filepath.Base(filepath.Dir(goModPath))

	log.WithFields(log.Fields{"from": parentDir, "to": copyDir}).Debug("Copying files")
	copy.Copy(parentDir, copyDir)

	logFiles(copyDir, "Copied Files")

	vendorDir(copyDir)

	updateVeracodeJson(mainFile, parentDir, copyDir)

	pkg(goModPath, tempWorkDir, parentDir, packageDate)

	logFiles(tempWorkDir, "Temporary workdir after packaging")
	logFiles(parentDir, "ParentDir")
	logFiles(parentDir+string(filepath.Separator)+"veracode", "Veracode Dir")
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
	ioutil.WriteFile(copyDir+string(filepath.Separator)+"veracode.json", json, 0644)
}

func pkg(goModPath string, tempWorkDir string, parentDir string, packageDate string) {
	goModDir := filepath.Dir(goModPath)
	log.Debug(goModDir)
	baseDir := filepath.Base(goModDir)
	if packageDate == "" {
		packageDate = time.Now().Format("-20060102150405")
	}
	zipFile := baseDir + packageDate + ".zip"
	log.WithFields(log.Fields{
		"baseDir": baseDir,
		"zipFile": tempWorkDir + string(filepath.Separator) + zipFile,
	}).Debug("Writing zip file")
	ZipWriter(tempWorkDir, tempWorkDir+string(filepath.Separator)+zipFile)

	veracodeDir := parentDir + string(filepath.Separator) + "veracode"
	os.Mkdir(veracodeDir, 0700)
	log.WithField("veracodeDir", veracodeDir).Debug("Created veracode dir for binaries")

	err := MoveFile(
		tempWorkDir+string(filepath.Separator)+zipFile,
		veracodeDir+string(filepath.Separator)+zipFile,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"from": tempWorkDir + string(filepath.Separator) + zipFile,
		"to":   veracodeDir + string(filepath.Separator) + zipFile,
	}).Debug("Rename zipfile")
}

func ZipWriter(baseFolder string, outputZipFile string) {
	// Get a Buffer to Write To
	outFile, err := os.Create(outputZipFile)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, baseFolder, "", filepath.Base(outputZipFile))

	if err != nil {
		log.Fatal(err)
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func addFiles(w *zip.Writer, basePath, baseInZip string, ignoreFile string) {
	// Open the Directory
	basePath = basePath + string(filepath.Separator)
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == ignoreFile {
			continue
		}

		log.Debug(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				log.Fatal(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				log.Fatal(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + string(filepath.Separator)
			log.Debug("Recursing and Adding SubDir: " + file.Name())
			log.Debug("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, baseInZip+file.Name()+string(filepath.Separator), "")
		}
	}
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		log.Fatal(err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		log.Fatal(err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
