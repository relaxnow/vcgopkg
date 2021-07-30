package main

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func ZipWriter(baseFolder string, outputZipFilePath string) error {
	// Create file.
	outFile, err := os.Create(outputZipFilePath)
	if outFile != nil {
		defer outFile.Close()
	}
	if err != nil {
		log.WithFields(log.Fields{
			"outputZipFile": outputZipFilePath,
			"err":           err,
		}).Error("Failed creating zip file")
		return err
	}

	// Create a new zip archive.
	zipWriter := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(zipWriter, baseFolder, "", filepath.Base(outputZipFilePath))

	if err != nil {
		return err
	}

	// Make sure to check the error on Close.
	return zipWriter.Close()
}

func addFiles(w *zip.Writer, basePath, baseInZip string, ignoreFile string) error {
	// Open the Directory
	basePath = basePath + string(filepath.Separator)
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Name() == ignoreFile {
			continue
		}

		log.Debug(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				log.Error(err)
				return err
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				log.Error(err)
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				log.Error(err)
				return err
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
