package pkg

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

func TestEnvConfig_Render(t *testing.T) {
	cfg := EnvConfig{
		Variables: []string{"PATH", "FOO", "BAR"},
	}

	r, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, EnvResult{}, r)

	result := r.(EnvResult)
	assert.Len(t, result.Env, 3)
	assert.Contains(t, result.Env, "PATH")
	assert.Contains(t, result.Env, "FOO")
	assert.Contains(t, result.Env, "BAR")
}

func TestEnvConfig_Render_None(t *testing.T) {
	cfg := EnvConfig{
		Variables: []string{},
	}

	r, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, EnvResult{}, r)

	result := r.(EnvResult)
	assert.Len(t, result.Env, 0)
}

func TestNewEnvResult(t *testing.T) {
	r := NewEnvResult()
	assert.Empty(t, r.Env)
}

func TestEnvResult_IsEmpty(t *testing.T) {
	r := NewEnvResult()
	assert.True(t, r.IsEmpty())

	r.Env["foo"] = "bar"
	assert.False(t, r.IsEmpty())
}

func TestEnvResult_Markdown(t *testing.T) {
	r := NewEnvResult()
	r.Env["FOO"] = "bar"
	r.Env["ABC"] = "123"

	data, err := r.Markdown()
	assert.NoError(t, err)

	expected := "**Environment**\n```\nABC=123\nFOO=bar\n```\n"
	assert.Equal(t, expected, string(data))
}

func TestEnvResult_Markdown_Empty(t *testing.T) {
	r := NewEnvResult()

	data, err := r.Markdown()
	assert.NoError(t, err)
	assert.Empty(t, data, string(data))
}

func TestEnvResult_Plaintext(t *testing.T) {
	r := NewEnvResult()
	r.Env["FOO"] = "bar"
	r.Env["ABC"] = "123"

	data, err := r.Plaintext()
	assert.NoError(t, err)

	expected := "Environment\n-----------\nABC=123\nFOO=bar\n"
	assert.Equal(t, expected, string(data))
}

func TestEnvResult_Plaintext_Empty(t *testing.T) {
	r := NewEnvResult()

	data, err := r.Plaintext()
	assert.NoError(t, err)
	assert.Empty(t, data, string(data))
}

func TestEnvResult_JSON(t *testing.T) {
	r := NewEnvResult()
	r.Env["FOO"] = "bar"
	r.Env["ABC"] = "123"

	data, err := r.JSON()
	assert.NoError(t, err)

	expected := `{"env":{"ABC":"123","FOO":"bar"}}`
	assert.Equal(t, expected, string(data))
}

func TestEnvResult_JSON_Empty(t *testing.T) {
	r := NewEnvResult()

	data, err := r.JSON()
	assert.NoError(t, err)
	assert.Empty(t, data, string(data))
}

func TestEnvResult_YAML(t *testing.T) {
	r := NewEnvResult()
	r.Env["FOO"] = "bar"
	r.Env["ABC"] = "123"

	data, err := r.YAML()
	assert.NoError(t, err)

	expected := heredoc.Doc(`
		ABC: "123"
		FOO: bar
	`)
	assert.Equal(t, expected, string(data))
}

func TestEnvResult_YAML_Empty(t *testing.T) {
	r := NewEnvResult()

	data, err := r.YAML()
	assert.NoError(t, err)
	assert.Empty(t, data, string(data))
}
