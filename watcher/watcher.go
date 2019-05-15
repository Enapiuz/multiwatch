package watcher

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Enapiuz/multiwatch/types"
	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
)

// Interface Watcher interface
type Interface interface {
	Run(chan bool)
	GetStatus() string
	GetErrors() string
}

// Watcher watch for directory changes and run commands
type Watcher struct {
	config   *types.DirectoryConfig
	watcher  *fsnotify.Watcher
	status   string
	errorLog string
}

//NewWatcher creates new watcher
func NewWatcher(config types.DirectoryConfig) *Watcher {
	watcher, _ := fsnotify.NewWatcher()
	newWatcher := &Watcher{
		status:   "⚪",
		config:   &config,
		watcher:  watcher,
		errorLog: "",
	}
	newWatcher.registerFiles()
	return newWatcher
}

// Run watcher
func (w *Watcher) Run(needReprint chan bool) {
	go func() {
		f := func() {
			w.status = "🔄"
			needReprint <- true
			result := w.runCommands()
			if result {
				w.status = "👍"
			} else {
				w.status = "🔴"
			}
			needReprint <- true
		}
		debounced := debounce.New(500 * time.Millisecond)
		debounced(f)
		for {
			select {
			// watch for events
			case event := <-w.watcher.Events:
				if event.Op.String() != "CHMOD" {
					debounced(f)
				}

			// watch for errors
			case err := <-w.watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()
}

func (w *Watcher) runCommands() bool {
	result := true
	w.errorLog = ""
	for _, command := range w.config.Commands {
		var outBuffer, errBuffer bytes.Buffer
		cmd := exec.Command("sh", "-c", command)
		cmd.Stdout = &outBuffer
		cmd.Stderr = &errBuffer
		err := cmd.Run()
		if err != nil {
			w.errorLog += fmt.Sprintf("[%s]:\n%s%s", command, outBuffer.String(), errBuffer.String())
			result = false
			if w.config.BreakOnFail {
				break
			}
		}
	}
	return result
}

// GetStatus returns text representation of current watcher's status
func (w *Watcher) GetStatus() string {
	return fmt.Sprintf("%s %s", w.status, w.config.Name)
}

// GetErrors returns errors from last watcher run
func (w *Watcher) GetErrors() string {
	return w.errorLog
}

func (w *Watcher) registerFiles() error {
	for _, dir := range w.config.Paths {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Fatalf("Directory '%s' in '%s' watcher does not exists", dir, w.config.Name)
		}

		if err := filepath.Walk(dir, w.watchDir(dir)); err != nil {
			log.Fatal("ERROR: ", err)
		}
	}
	return nil
}

func (w *Watcher) watchDir(baseDir string) func(path string, fi os.FileInfo, err error) error {
	return func(path string, fi os.FileInfo, err error) error {
		if fi.Mode().IsDir() {
			// check absolute path
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			// check ignore prefixes
			absBaseDir, err := filepath.Abs(baseDir)
			if err != nil {
				return err
			}

			for _, ignorePrefix := range w.config.IgnorePrefixes {
				targetPath := fmt.Sprintf("%s/%s", absBaseDir, ignorePrefix)
				if strings.HasPrefix(absPath, targetPath) {
					return nil
				}
			}
			return w.watcher.Add(absPath)
		}
		return nil
	}
}
