package watcher

import (
	"fmt"
	"testing"

	"github.com/Enapiuz/multiwatch/types"
	"github.com/stretchr/testify/assert"
)

func TestWatcher_GetStatus(t *testing.T) {
	conf := types.DirectoryConfig{
		Name:     "testwatcher",
		Paths:    []string{"../watcher"},
		Commands: []string{""},
	}
	watcher := NewWatcher(conf)
	status := watcher.GetStatus()
	expected := fmt.Sprintf("%s %s", "⚪", "testwatcher")
	assert.Equal(t, expected, status, "wrong default status")
}
