package main

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"github.com/relaxnow/go-detailedreport-to-csv/detailedreport"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"os"
)


func main() {
	xmlFilePath := os.Args[1]

	// Open our xmlFile
	xmlFile, err := os.Open(xmlFilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Successfully Opened " + xmlFilePath)
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)


	// we initialize our Users array
	var report detailedreport.DetailedReport

	reader := bytes.NewReader(byteValue)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&report)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//spew.Dump(report);
	rows := [][]string{{"Flaw ID", "CWE ID", "CWE Name", "Module", "Source File", "Severity"}}
	for i := 0; i < len(report.Severities); i++ {
		var severity = report.Severities[i]
		for j := 0; j < len(severity.Categories); j++ {
			var category = severity.Categories[j];
			for k := 0; k < len(category.CWES); k++ {
				var CWE = category.CWES[k];
				for l := 0; l < len(CWE.StaticFlaws.Flaws); l++ {
					var flaw = CWE.StaticFlaws.Flaws[l];

					rows = append(rows, []string{
						flaw.ID,
						CWE.ID,
						CWE.Name,
						flaw.Module,
						flaw.SourceFile,
						severity.Level})
				}
			}
		}
	}

	//spew.Dump(rows)

	csvfile,err := os.Create("output.csv")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	csvwriter := csv.NewWriter(csvfile)
	csvwriter.Comma = ';'

	for _, row := range rows {
		_ = csvwriter.Write(row)
	}

	csvwriter.Flush()

	csvfile.Close()
	fmt.Println("Wrote output.csv")
	fmt.Println("All done")
}
