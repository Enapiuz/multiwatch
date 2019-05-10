package printer

import (
	"fmt"
	"github.com/Enapiuz/multiwatch/watcher"
	"os"
	"os/exec"
	"runtime"
)

type Printer struct {
	watchers []*watcher.Watcher
	clear    map[string]func()
}

func NewPrinter() *Printer {
	clear := make(map[string]func())
	clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	return &Printer{clear: clear}
}

func (p *Printer) RegisterWatchers(watchers []*watcher.Watcher) {
	p.watchers = watchers
}

func (p *Printer) Start(needReprint chan bool) {
	go func() {
		for {
			select {
			case _ = <-needReprint:
				p.printWatchers()
			}
		}
	}()
}

func (p *Printer) printWatchers() {
	p.callClear()
	for _, localWatcher := range p.watchers {
		fmt.Println(localWatcher.GetStatus())
	}
}

func (p *Printer) callClear() {
	value, ok := p.clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                            //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
