package main

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/env-io/dev/core"
)

var (
	runCommand = &core.Command{
		Name: "run",
		Info: "start watching any changes in directory and rebuild it.",
		Usage: `
dev run
	Run command will watching any changes in directory of go project,
	it will recompile and restart the application binary.`,
	}
)

func init() {
	runCommand.Run = actionRun
}

// actionRun, perform to scan directory and get ready to watching.
func actionRun(_ *core.Command, args []string) int {
	gps := getGoPath()
	if len(gps) == 0 {
		core.LogError("$GOPATH not found, Please set $GOPATH in your environment variables.")
		os.Exit(2)
	}

	exit := make(chan bool)
	cwd, _ := os.Getwd()
	appName := path.Base(cwd)

	core.Log.Print("")
	if len(args) > 0 {
		called := args[0]

		if called != "" {
			appName = called
		}
	}

	var paths []string
	var files []string

	readDirectory(cwd, &paths)
	core.Watch(appName, paths, files)
	core.Build()

	for {
		select {
		case <-exit:
			runtime.Goexit()
		}
	}
}

// readDirectory binds paths with list of existing directory.
func readDirectory(directory string, paths *[]string) {
	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}

	useDirectory := false
	for _, fileInfo := range fileInfos {
		if strings.HasSuffix(fileInfo.Name(), "docs") {
			continue
		}

		if fileInfo.IsDir() == true && fileInfo.Name()[0] != '.' {
			readDirectory(directory+"/"+fileInfo.Name(), paths)
			continue
		}

		if useDirectory == true {
			continue
		}

		if path.Ext(fileInfo.Name()) == ".go" {
			*paths = append(*paths, directory)
			useDirectory = true
		}
	}

	return
}

// getGoPath returns list of go path on system.
func getGoPath() (p []string) {
	gopath := os.Getenv("GOPATH")
	p = strings.Split(gopath, ":")

	return
}
