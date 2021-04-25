package main

import (
	"io"
	"os"
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
