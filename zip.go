package main

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func ZipWriter(baseFolder string, outputZipFilePath string) {
	// Create file.
	outFile, err := os.Create(outputZipFilePath)
	if outFile != nil {
		defer outFile.Close()
	}
	if err != nil {
		log.WithFields(log.Fields{
			"outputZipFile": outputZipFilePath,
			"err":           err,
		}).Fatal("Failed creating zip file")
	}

	// Create a new zip archive.
	zipWriter := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(zipWriter, baseFolder, "", filepath.Base(outputZipFilePath))

	if err != nil {
		log.Fatal(err)
	}

	// Make sure to check the error on Close.
	err = zipWriter.Close()
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
