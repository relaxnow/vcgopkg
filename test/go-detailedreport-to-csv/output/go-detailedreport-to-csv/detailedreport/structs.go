package detailedreport

type DetailedReport struct {
	Severities []Severity		`xml:"severity"`
}
type Severity struct {
	Level  		string		`xml:"level,attr"`
	Categories	[]Category	`xml:"category"`
}
type Category struct {
	CWES	[]CWE	`xml:"cwe"`
}
type CWE struct {
	ID 			string		`xml:"cweid,attr"`
	Name 		string		`xml:"cwename,attr"`
	StaticFlaws StaticFlaws	`xml:"staticflaws"`
	Line		string		`xml:"line,attr"`
}
type StaticFlaws struct {
	Flaws	[]Flaw	`xml:"flaw"`
}
type Flaw struct {
	ID				string	`xml:"issueid,attr"`
	SourceFile		string	`xml:"sourcefile,attr"`
	SourceFilePath	string	`xml:"sourcefilepath,attr"`
	Module			string	`xml:"module,attr"`
}
