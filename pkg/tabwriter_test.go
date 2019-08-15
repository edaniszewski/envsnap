package pkg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTabWriter(t *testing.T) {
	tw := NewTabWriter(&bytes.Buffer{})
	assert.NotNil(t, tw)
}
