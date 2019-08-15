package pkg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWarnings(t *testing.T) {
	w := NewWarnings()
	assert.Len(t, w.Warnings, 0)
}

func TestWarnings_Add(t *testing.T) {
	w := NewWarnings()
	assert.Len(t, w.Warnings, 0)

	w.Add("src1", "msg: %d", 1)
	assert.Len(t, w.Warnings, 1)
}

func TestWarnings_Clear(t *testing.T) {
	w := NewWarnings()
	w.Warnings["foo"] = []string{"bar", "baz"}
	w.Warnings["abc"] = []string{"123"}
	assert.Len(t, w.Warnings, 2)

	w.Clear()
	assert.Len(t, w.Warnings, 0)
}

func TestWarnings_Clear_Empty(t *testing.T) {
	w := NewWarnings()
	assert.Len(t, w.Warnings, 0)

	w.Clear()
	assert.Len(t, w.Warnings, 0)
}

func TestWarnings_AddMultiple(t *testing.T) {
	w := NewWarnings()
	assert.Len(t, w.Warnings, 0)

	// first message
	w.Add("src1", "msg: %d", 1)
	assert.Len(t, w.Warnings, 1)
	assert.Len(t, w.Warnings["src1"], 1)

	// second message, same source
	w.Add("src1", "msg: %d", 2)
	assert.Len(t, w.Warnings, 1)
	assert.Len(t, w.Warnings["src1"], 2)

	// duplicate message
	w.Add("src1", "msg: %d", 2)
	assert.Len(t, w.Warnings, 1)
	assert.Len(t, w.Warnings["src1"], 3)
}

func TestWarnings_HasWarnings(t *testing.T) {
	w := NewWarnings()
	assert.False(t, w.HasWarnings())

	w.Warnings["foo"] = []string{"bar"}
	assert.True(t, w.HasWarnings())
}

func TestWarnings_PrintNoWarnings(t *testing.T) {
	out := bytes.Buffer{}
	w := NewWarnings()

	w.Print(&out)
	assert.Equal(t, "", out.String())
}

func TestWarnings_Print(t *testing.T) {
	out := bytes.Buffer{}
	w := NewWarnings()
	w.Warnings["foo"] = []string{"bar"}
	w.Warnings["abc"] = []string{"123"}

	expected := `------------------------------
warnings: 2

[abc]   123
[foo]   bar

`
	w.Print(&out)
	assert.Equal(t, expected, out.String())
}
