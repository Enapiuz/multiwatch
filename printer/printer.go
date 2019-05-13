package printer

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Enapiuz/multiwatch/watcher"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

// Printer watch for needReprint events and print all workers' statuses
type Printer struct {
	watchers []watcher.Interface
	clear    map[string]func()
}

// NewPrinter makes new printer object
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

// RegisterWatchers gets watchers to print when needReprint is toggled
func (p *Printer) RegisterWatchers(watchers []watcher.Interface) {
	p.watchers = watchers
}

// Start watching for reprint event
func (p *Printer) Start(needReprint chan bool) {
	go func() {
		for {
			select {
			case _ = <-needReprint:
				p.printWatchers()
				p.padToTop()
			}
		}
	}()
}

func (p *Printer) printWatchers() {
	p.callClear()
	var statuses strings.Builder
	for idx, localWatcher := range p.watchers {
		if errorText := localWatcher.GetErrors(); errorText != "" {
			fmt.Println(errorText)
		}
		var appendText strings.Builder
		if idx == 0 {
			appendText.WriteString(fmt.Sprint(localWatcher.GetStatus()))
			width, err := terminal.Width()
			if err == nil {
				currentTime := fmt.Sprintf("%s\n", time.Now().Format("15:04:05"))
				toRepeat := int(width) - utf8.RuneCountInString(currentTime) - utf8.RuneCountInString(appendText.String())
				if toRepeat < 0 {
					toRepeat = 0
				}
				appendText.WriteString(
					strings.Repeat(
						" ",
						toRepeat,
					),
				)
				appendText.WriteString(currentTime)
			}
		} else {
			appendText.WriteString(fmt.Sprintln(localWatcher.GetStatus()))
		}
		statuses.WriteString(appendText.String())
	}
	fmt.Println(strings.Trim(statuses.String(), "\n"))
}

func (p *Printer) callClear() {
	value, ok := p.clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func (p *Printer) padToTop() {
	y, err := terminal.Height()
	if err == nil {
		toRepeat := int(y) - len(p.watchers) - 1
		if toRepeat < 0 {
			toRepeat = 0
		}
		fmt.Print(strings.Repeat("\n", toRepeat))
	}
	x, err := terminal.Width()
	if err == nil {
		toRepeat := int(x) - 1
		if toRepeat < 0 {
			toRepeat = 0
		}
		fmt.Print(strings.Repeat(" ", toRepeat))
	}
}
