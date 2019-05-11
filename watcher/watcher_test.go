package watcher

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWatcher_GetStatus(t *testing.T) {
	watcher := NewWatcher("testwatcher", []string{"../watcher"}, []string{""})
	status := watcher.GetStatus()
	expected := fmt.Sprintf("%s %s", "⚪", "testwatcher")
	assert.Equal(t, expected, status, "wrong default status")
}
