package printer

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/Enapiuz/multiwatch/watcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type WatcherMock struct {
	mock.Mock
}

func (w *WatcherMock) Run(chan bool) {
	panic("implement me")
}

func (w *WatcherMock) GetStatus() string {
	return "status_log"
}

func (w *WatcherMock) GetErrors() string {
	return "error_log"
}

func TestNewPrinter(t *testing.T) {
	samplePrinter := NewPrinter()
	assert.Len(t, samplePrinter.watchers, 0)
}

func TestPrinter_RegisterWatchers(t *testing.T) {
	testPrinter := NewPrinter()
	dumbWatcher := &WatcherMock{}
	testPrinter.RegisterWatchers([]watcher.Interface{dumbWatcher})
	assert.Len(t, testPrinter.watchers, 1)
}

func TestPrinter_Start(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	testPrinter := NewPrinter()
	dumbWatcher := &WatcherMock{}
	testPrinter.RegisterWatchers([]watcher.Interface{dumbWatcher})
	needReprint := make(chan bool)
	testPrinter.Start(needReprint)
	needReprint <- true

	time.Sleep(1 * time.Second)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	assert.Contains(t, string(out), "status_log")
	assert.Contains(t, string(out), "error_log")
}
