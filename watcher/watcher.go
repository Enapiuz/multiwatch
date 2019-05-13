package watcher

import (
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

type Watcher struct {
	config  *types.DirectoryConfig
	watcher *fsnotify.Watcher
	status  string
}

func NewWatcher(config types.DirectoryConfig) *Watcher {
	watcher, _ := fsnotify.NewWatcher()
	newWatcher := &Watcher{
		status:  "⚪",
		config:  &config,
		watcher: watcher,
	}
	newWatcher.registerFiles()
	return newWatcher
}

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
			case _ = <-w.watcher.Events:
				debounced(f)

			// watch for errors
			case err := <-w.watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()
}

func (w *Watcher) runCommands() bool {
	result := true
	for _, command := range w.config.Commands {
		cmd := exec.Command("sh", "-c", command)
		err := cmd.Run()
		if err != nil {
			result = false
		}
	}
	return result
}

func (w *Watcher) GetStatus() string {
	return fmt.Sprintf("%s %s", w.status, w.config.Name)
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
