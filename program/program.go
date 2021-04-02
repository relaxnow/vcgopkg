package program

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/relaxnow/vcgopkg/files"
)

type PackagableProgram interface {
	CopyToTempDir()
	Vendor()
	Zip()
}

type program struct {
	MainFilePath string
	RootPath     string
	Files        files.FeatureFiles
}

type workspaceProgram struct {
	Program program
	GoPath  string
}

func (p *workspaceProgram) CopyToTempDir() {}
func (p *workspaceProgram) Vendor()        {}
func (p *workspaceProgram) Zip()           {}

type moduleProgram struct {
	Program    program
	ModuleName string
}

func (p *moduleProgram) CopyToTempDir() {}
func (p *moduleProgram) Vendor()        {}
func (p *moduleProgram) Zip()           {}

func GetProgramFromMainFilePath(mainFilePath string, inputPathDir string) *workspaceProgram {
	gopath := detectMatchingGoPath(mainFilePath)
	if gopath == "" {
		return &workspaceProgram{
			Program: program{
				MainFilePath: mainFilePath,
				RootPath:     inputPathDir,
				Files:        files.FeatureFiles{},
			},
			GoPath: gopath,
		}
	}
	return nil
}

func detectMatchingGoPath(filePath string) string {
	gopaths := []string{
		os.Getenv("GOPATH"),
		getGoPathFromGoEnv(),
	}

	for _, gopath := range gopaths {
		if gopath == "" {
			// TODO warn
			return ""
		}

		if strings.HasPrefix(filePath, gopath) {
			return gopath
		}
	}
	return ""
}

// Detecting a program root in workspace mode is actually a bit difficult.
// src/github.com/user/program/cmd/server/main.go
func detectProgramRoot(mainFilePath string, gopath string) string {
	return ""
}

func getGoPathFromGoEnv() string {
	goCmdPath, _ := exec.LookPath("go")
	if goCmdPath == "" {
		return ""
	}

	fmt.Println("Go is at: " + goCmdPath)
	cmd := exec.Command("go", "env")
	cmdOut, _ := cmd.Output()
	variables := strings.Split(string(cmdOut), "\n")

	for _, variable := range variables {
		if strings.HasPrefix(variable, "GOPATH=") {
			return strings.TrimPrefix(variable, "GOPATH=")
		}
	}
	return ""
}

// func hasGoMod(path) {
// 	mod, err := modfile.Parse(path)
// }

// func isWorkspaceModeWithGoEnv(path) {
// 	// if !GOPATH
// 	//    loadEnvFromGoEnv()
// 	// return isWorkspaceMode(path)
// }

// func isWorkspaceMode(path) {
// 	// workspaceRoot := getWorkspaceRoot()
// 	// if workspaceRoot == nil {
// 	//    return false
// 	// }
// 	// return path.startsWith(path)
// }

// func getWorkspaceRoot() {
// 	// if isWorkspace(GOPATH)
// 	//   return GOPATH
// 	// if isWorkspace(GOROOT)
// 	//    return GOROOT
// 	// return nil
// }

// func isWorkspace(path) {
// 	// if !path.Stat().isDir()
// 	//   return false
// 	// if !path + '/src'.isDir()
// 	//   return false
// 	// return true
// }
