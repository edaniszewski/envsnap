package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecConfig_Render(t *testing.T) {
	cfg := ExecConfig{
		Run: []string{
			`echo "testing"`,
		},
	}

	r, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, ExecResult{}, r)

	result := r.(ExecResult)
	assert.Len(t, result.Exec, 1)
	assert.Equal(t, "\"testing\"\n", result.Exec[`echo "testing"`])
}

func TestExecConfig_Render_Err(t *testing.T) {
	defer cliWarnings.Clear()

	cfg := ExecConfig{
		Run: []string{
			`ls xyz`,
		},
	}

	r, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, ExecResult{}, r)

	result := r.(ExecResult)
	assert.Len(t, result.Exec, 1)
	assert.Equal(t, "", result.Exec[`ls xyz`])

	assert.Contains(t, cliWarnings.Warnings, "exec.run")
	assert.Len(t, cliWarnings.Warnings["exec.run"], 1)
}

func TestExecConfig_Render_None(t *testing.T) {
	cfg := ExecConfig{
		Run: []string{},
	}

	r, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, ExecResult{}, r)

	result := r.(ExecResult)
	assert.Len(t, result.Exec, 0)
}

func TestNewExecResult(t *testing.T) {
	r := NewExecResult()
	assert.Empty(t, r.Exec)
}

func TestExecResult_IsEmpty(t *testing.T) {
	r := NewExecResult()
	assert.True(t, r.IsEmpty())

	r.Exec["foo"] = "bar"
	assert.False(t, r.IsEmpty())
}

func TestExecResult_Markdown(t *testing.T) {
	r := NewExecResult()
	r.Exec["echo hello"] = "hello"

	data, err := r.Markdown()
	assert.NoError(t, err)

	expected := "**Exec**\n- `echo hello`\n  ```\n  hello\n  ```\n"
	assert.Equal(t, expected, string(data))
}

func TestExecResult_Markdown_Empty(t *testing.T) {
	r := NewExecResult()

	data, err := r.Markdown()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

func TestExecResult_Plaintext(t *testing.T) {
	r := NewExecResult()
	r.Exec["echo hello"] = "hello"

	data, err := r.Plaintext()
	assert.NoError(t, err)

	expected := "Exec\n----\n$ echo hello\n  hello\n"
	assert.Equal(t, expected, string(data))
}

func TestExecResult_Plaintext_Empty(t *testing.T) {
	r := NewExecResult()

	data, err := r.Plaintext()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

func TestExecResult_JSON(t *testing.T) {
	r := NewExecResult()
	r.Exec["echo hello"] = "hello"

	data, err := r.JSON()
	assert.NoError(t, err)

	expected := `{"exec":{"echo hello":"hello"}}`
	assert.Equal(t, expected, string(data))
}

func TestExecResult_JSON_Empty(t *testing.T) {
	r := NewExecResult()

	data, err := r.JSON()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

func TestExecResult_YAML(t *testing.T) {
	r := NewExecResult()
	r.Exec["echo hello"] = "hello"

	data, err := r.YAML()
	assert.NoError(t, err)

	expected := "exec:\n  echo hello: hello\n"
	assert.Equal(t, expected, string(data))
}

func TestExecResult_YAML_Empty(t *testing.T) {
	r := NewExecResult()

	data, err := r.YAML()
	assert.NoError(t, err)
	assert.Empty(t, data)
}
