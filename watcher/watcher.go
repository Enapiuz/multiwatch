package watcher

import (
	"fmt"
	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Watcher struct {
	name     string
	watcher  *fsnotify.Watcher
	commands []string
	status   string
}

func NewWatcher(name string, dirs []string, commands []string) *Watcher {
	watcher, _ := fsnotify.NewWatcher()
	for _, dir := range dirs {
		if err := filepath.Walk(dir, watchDir(watcher)); err != nil {
			fmt.Println("ERROR", err)
		}
	}
	return &Watcher{watcher: watcher, commands: commands, name: name, status: "⚪"}
}

func (w *Watcher) Run(needReprint chan bool) {
	log.Printf("watcher %s has been started", w.name)
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
	for _, command := range w.commands {
		cmd := exec.Command("sh", "-c", command)
		err := cmd.Run()
		if err != nil {
			result = false
		}
	}
	return result
}

func (w *Watcher) GetStatus() string {
	//return w.name
	return fmt.Sprintf("%s %s", w.status, w.name)
}

func watchDir(watcher *fsnotify.Watcher) func(path string, fi os.FileInfo, err error) error {
	return func(path string, fi os.FileInfo, err error) error {
		if fi.Mode().IsDir() {
			return watcher.Add(path)
		}

		return nil
	}
}
