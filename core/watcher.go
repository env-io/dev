package core

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/howeyc/fsnotify"
)

var (
	// appName, packages name, get from directory name.
	appName string

	// cmd is external command.
	cmd *exec.Cmd

	// state is mutual exclusion lock,
	state sync.Mutex

	// watching, list of extention file that need to watch.
	watching []string

	// lastBuild, time last build performed.
	lastBuild time.Time

	// isStarted, if true then application is on running.
	isStarted = make(chan bool)

	// runTime, slices of runnable application unix times.
	runTime = make(map[string]int64)
)

type DocVal string

func (d *DocVal) String() string {
	return fmt.Sprint(*d)
}

func (d *DocVal) Set(value string) error {
	*d = DocVal(value)
	return nil
}

type ListOpts []string

func (opts *ListOpts) String() string {
	return fmt.Sprint(*opts)
}

func (opts *ListOpts) Set(value string) error {
	*opts = append(*opts, value)
	return nil
}

// Watch performs initializing fsnotify to watch files on current directory.
// and trigger to rebuild and restart the packages after some changes has been made
// on files which that we watch.
func Watch(appname string, paths []string, files []string) {
	appName = appname
	watching = append(files, ".go")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		LogError("Failed to create new Watcher")
		LogError(err.Error())
		os.Exit(2)
	}

	go func() {
		for {
			select {
			case e := <-watcher.Event:
				isbuild := true
				if !isWatched(e.Name) {
					continue
				}

				if lastBuild.Add(1 * time.Second).After(time.Now()) {
					continue
				}

				lastBuild = time.Now()
				mt := lastModified(e.Name)
				if t := runTime[e.Name]; mt == t {
					LogInfo(fmt.Sprintf("Skipped # %s #\n", e.String()))
					isbuild = false
				}

				runTime[e.Name] = mt
				if isbuild {
					LogInfo(e.String())
					go Build()
				}
			case err := <-watcher.Error:
				LogError(err.Error())
			}
		}
	}()

	for i, path := range paths {
		bar := progress(i+1, len(paths), 120)
		os.Stdout.Write([]byte(bar + "\r"))

		err = watcher.Watch(path)
		if err != nil {
			LogError("Failed to watch directory.")
			LogError(err.Error())
			os.Exit(2)
		}

		os.Stdout.Sync()
	}

	os.Stdout.Write([]byte("\n"))
}

// Build performs executing go build of packages.
func Build() {
	var err error

	state.Lock()
	defer state.Unlock()

	path, _ := os.Getwd()
	os.Chdir(path)

	// using script from go-fast-build
	// github.com/kovetskiy/go-fast
	cmd := exec.Command("/bin/sh", "-c", "export GOBIN=$(pwd); exec go install -gcflags \"-trimpath $GOPATH/src\" \"$@\";")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		LogError("Failed to build.")
		LogError(err.Error())
	}

	restart(appName)
}

func runTest() {
	LogInfo("Running vet, lint and test ...")
	cmd := exec.Command("/bin/sh", "-c", "go test -cover;")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		LogError("Failed to run test.")
		LogError(err.Error())
	}
}

// restart performs restarting application binary.
func restart(app string) {
	kill()

	go start(app)
}

// kill performs killing current running application.
func kill() {
	defer func() {
		if e := recover(); e != nil {
			LogError(fmt.Sprintf("Kill.recover -> %s", e))
		}
	}()

	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
	}
}

// start performs to start the application binary.
func start(app string) {
	LogInfo(ColorRed(fmt.Sprintf("Rebuild %s", app)))

	if strings.Index(app, "./") == -1 {
		app = "./" + app
	}

	cmd = exec.Command(app)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go cmd.Run()
	LogInfo(ColorGreen(fmt.Sprintf("Starting %s", app)))
	isStarted <- true
}

// lastModified returns unix timestamp of file last modified.
func lastModified(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		LogError(fmt.Sprintf("Cannot find file [ %s ]\n", err))
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		LogError(fmt.Sprintf("Failed to get information from file [ %s ]\n", err))
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

// isWatched returns if ext files was on watching.
func isWatched(fileName string) bool {
	for _, s := range watching {
		if strings.HasSuffix(fileName, s) {
			return true
		}
	}

	return false
}

// progress returns a string as progress bar from scaning directory.
func progress(current, total, cols int) string {
	prefix := strconv.Itoa(current) + " / " + strconv.Itoa(total)
	barStart := " ["
	barEnd := "] "

	barSize := cols - len(prefix+barStart+barEnd)
	amount := int(float32(current) / (float32(total) / float32(barSize)))
	remain := barSize - amount

	bar := strings.Repeat("#", amount) + strings.Repeat(" ", remain)
	return ColorDim("Scanning Directory \033[1m" + prefix + "\033[0m" + barStart + bar + barEnd)
}
