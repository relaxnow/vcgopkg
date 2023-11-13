package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

func CreateEmptyVeracodeJsonFileIfNotExists(file string) error {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		os.WriteFile(file, []byte("{}"), os.FileMode(int(0666)))
		return nil
	} else {
		return err
	}
}

type VeracodeJson struct {
	MainRoot    string   `json:"MainRoot,omitempty"`
	MainPkgName string   `json:"MainPkgName,omitempty"`
	FirstParty  []string `json:"FirstParty,omitempty"`
}

func NewVeracodeJson() VeracodeJson {
	return VeracodeJson{}
}

func NewVeracodeJsonFromFile(file string) (VeracodeJson, error) {
	veracodeJsonFile, err := os.Open(file)

	if err != nil {
		return VeracodeJson{}, err
	}

	defer veracodeJsonFile.Close()

	byteValue, err := io.ReadAll(veracodeJsonFile)

	if err != nil {
		return VeracodeJson{}, err
	}

	var veracodeJson VeracodeJson
	err = json.Unmarshal(byteValue, &veracodeJson)

	if err != nil {
		return VeracodeJson{}, err
	}

	return veracodeJson, nil
}

func (veracodeJson VeracodeJson) WriteToFile(file string) error {
	contents, err := json.Marshal(veracodeJson)

	if err != nil {
		return err
	}

	err = os.WriteFile(file, contents, 0644)

	if err != nil {
		return err
	}

	return err
}

type VeracodeJsonFile struct {
	File         string
	VeracodeJson VeracodeJson
}

func NewVeracodeJsonFile(file string) (VeracodeJsonFile, error) {
	veracodeJsonFile, err := NewVeracodeJsonFromFile(file)

	if err != nil {
		return VeracodeJsonFile{File: file}, err
	}

	return VeracodeJsonFile{File: file, VeracodeJson: veracodeJsonFile}, nil
}

func (veracodeJsonFile VeracodeJsonFile) WriteToFile() error {
	return veracodeJsonFile.VeracodeJson.WriteToFile(veracodeJsonFile.File)
}
