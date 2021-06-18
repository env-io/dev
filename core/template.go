package core

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type StubTemplate struct {
	AppPath           string
	ProjectPath       string
	PackagePath       string
	PackageName       string
	ModuleName        string
	ModelName         string
	ModelNameSingular string
	ModelNamePlural   string
	TableName         string
}

func GetDirName(currentPath string) string {
	LogInfo(currentPath)

	var x = strings.Split(currentPath, "/")

	return x[len(x)-1]
}

func GetPackagePath(currentPath string) string {
	gp := os.Getenv("GOPATH")
	if gp == "" {
		LogError("you should set GOPATH in the env")
		os.Exit(2)
	}

	appPath := ""
	haspath := false
	for _, wg := range filepath.SplitList(gp) {
		wg, _ = filepath.EvalSymlinks(path.Join(wg, "src"))

		if filepath.HasPrefix(strings.ToLower(currentPath), strings.ToLower(wg)) {
			haspath = true
			appPath = wg
			break
		}
	}

	if !haspath {
	}

	return strings.Join(strings.Split(currentPath[len(appPath)+1:], string(filepath.Separator)), "/")
}

func GetGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		LogError("$GOPATH not found, Please set $GOPATH in your environment variables.")
		os.Exit(2)
	}

	return gopath
}

func FormatSourceCode(filename string) {
	cmd := exec.Command("gofmt", "-w", filename)
	if err := cmd.Run(); err != nil {
		LogError("gofmt err:")
		LogError(err.Error())
	}
}
