package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Move file instead of renaming as renaming does not work across devices
// This is necessary as the temporary directory we use for packaging
// may be on a different disk.
func MoveFile(sourcePath, destPath string) error {
	// Open input file
	inputFile, err := os.Open(sourcePath)
	if inputFile != nil {
		defer inputFile.Close()
	}
	if err != nil {
		return err
	}

	// Open output file
	outputFile, err := os.Create(destPath)
	if outputFile != nil {
		defer outputFile.Close()
	}
	if err != nil {
		return err
	}

	// Copy file
	_, err = io.Copy(outputFile, inputFile)

	return err
}

func deleteFilesWithoutGoExtension(root string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if the file has a ".go" extension
		if !strings.HasSuffix(info.Name(), ".go") {
			log.Debug("deleteFilesWithoutGoExtension: Deleting: " + path)
			if err := os.Remove(path); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
