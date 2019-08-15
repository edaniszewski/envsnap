package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	app := NewApp()

	assert.Equal(t, "envsnap", app.Name)
	assert.Equal(t, Version, app.Version)
	assert.Len(t, app.Commands, 2)
}
