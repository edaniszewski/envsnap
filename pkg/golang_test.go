package pkg

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

func TestGolangConfig_Render(t *testing.T) {
	cfg := GolangConfig{
		Core: []string{"version"},
	}

	out, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, GolangResult{}, out)

	res := out.(GolangResult)
	assert.NotEmpty(t, res.Version)
	assert.Empty(t, res.Goroot)
	assert.Empty(t, res.Gopath)
}

func TestGolangConfig_Render2(t *testing.T) {
	cfg := GolangConfig{
		Core: []string{"goroot"},
	}

	out, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, GolangResult{}, out)

	res := out.(GolangResult)
	assert.Empty(t, res.Version)
	assert.NotEmpty(t, res.Goroot)
	assert.Empty(t, res.Gopath)
}

func TestGolangConfig_Render3(t *testing.T) {
	cfg := GolangConfig{
		Core: []string{"gopath"},
	}

	out, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, GolangResult{}, out)

	res := out.(GolangResult)
	assert.Empty(t, res.Version)
	assert.Empty(t, res.Goroot)
	assert.NotEmpty(t, res.Gopath)
}

func TestGolangConfig_Render_Err(t *testing.T) {
	cfg := GolangConfig{
		Core: []string{"not-an-option"},
	}

	_, err := cfg.Render()
	assert.Error(t, err)
}

func TestGolangConfig_Render_Empty(t *testing.T) {
	cfg := GolangConfig{}

	out, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, GolangResult{}, out)

	res := out.(GolangResult)
	assert.True(t, res.IsEmpty())
}

func TestNewGolangResult(t *testing.T) {
	r := NewGolangResult()
	assert.Empty(t, r.Version)
	assert.Empty(t, r.Goroot)
	assert.Empty(t, r.Gopath)
}

func TestGolangResult_IsEmpty(t *testing.T) {
	r := NewGolangResult()
	assert.True(t, r.IsEmpty())
}

func TestGolangResult_IsEmpty2(t *testing.T) {
	r := NewGolangResult()
	r.Version = "1"
	assert.False(t, r.IsEmpty())
}

func TestGolangResult_IsEmpty3(t *testing.T) {
	r := NewGolangResult()
	r.Goroot = "/path"
	assert.False(t, r.IsEmpty())
}

func TestGolangResult_IsEmpty4(t *testing.T) {
	r := NewGolangResult()
	r.Gopath = "/path"
	assert.False(t, r.IsEmpty())
}

func TestGolangResult_Markdown(t *testing.T) {
	r := NewGolangResult()
	r.Version = "1.13"
	r.Goroot = "/go/root"
	r.Gopath = "/go/path"

	data, err := r.Markdown()
	assert.NoError(t, err)

	expected := "**Golang**\n- _version_: 1.13\n- _goroot_: /go/root\n- _gopath_: /go/path\n"
	assert.Equal(t, expected, string(data))
}

func TestGolangResult_Markdown_Empty(t *testing.T) {
	r := NewGolangResult()

	data, err := r.Markdown()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

func TestGolangResult_Plaintext(t *testing.T) {
	r := NewGolangResult()
	r.Version = "1.13"
	r.Goroot = "/go/root"
	r.Gopath = "/go/path"

	data, err := r.Plaintext()
	assert.NoError(t, err)

	expected := "Golang\n------\nversion:  1.13\ngoroot:   /go/root\ngopath:   /go/path\n"
	assert.Equal(t, expected, string(data))
}

func TestGolangResult_Plaintext_Empty(t *testing.T) {
	r := NewGolangResult()

	data, err := r.Plaintext()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

func TestGolangResult_JSON(t *testing.T) {
	r := NewGolangResult()
	r.Version = "1.13"
	r.Goroot = "/go/root"
	r.Gopath = "/go/path"

	data, err := r.JSON()
	assert.NoError(t, err)

	expected := `{"version":"1.13","goroot":"/go/root","gopath":"/go/path"}`
	assert.Equal(t, expected, string(data))
}

func TestGolangResult_JSON_Empty(t *testing.T) {
	r := NewGolangResult()

	data, err := r.JSON()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

func TestGolangResult_YAML(t *testing.T) {
	r := NewGolangResult()
	r.Version = "1.13"
	r.Goroot = "/go/root"
	r.Gopath = "/go/path"

	data, err := r.YAML()
	assert.NoError(t, err)

	expected := heredoc.Doc(`
		version: "1.13"
		goroot: /go/root
		gopath: /go/path
	`)
	assert.Equal(t, expected, string(data))
}

func TestGolangResult_YAML_Empty(t *testing.T) {
	r := NewGolangResult()

	data, err := r.YAML()
	assert.NoError(t, err)
	assert.Empty(t, data)
}
