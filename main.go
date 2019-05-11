package main

import (
	"github.com/BurntSushi/toml"
	"github.com/Enapiuz/multiwatch/printer"
	"github.com/Enapiuz/multiwatch/watcher"
	"log"
)

type DirectoryConfig struct {
	Name     string
	Paths    []string
	Commands []string
}

type Config struct {
	Delay int32
	Watch []DirectoryConfig
}

func main() {
	var config Config
	var watchers = make([]*watcher.Watcher, 0)
	needReprint := make(chan bool)
	_, err := toml.DecodeFile("multiwatch.toml", &config)
	if err != nil {
		log.Fatal(err)
	}
	for _, watchConfig := range config.Watch {
		dirWatcher := watcher.NewWatcher(watchConfig.Name, watchConfig.Paths, watchConfig.Commands)
		dirWatcher.Run(needReprint)
		watchers = append(watchers, dirWatcher)
	}

	var statusPrinter = printer.NewPrinter()
	statusPrinter.RegisterWatchers(watchers)
	statusPrinter.Start(needReprint)
	needReprint <- true

	done := make(chan bool)
	<-done
}
