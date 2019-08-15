package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// GolangConfig defines the configuration for the "go" source.
type GolangConfig struct {
	Core []string `yaml:"core,omitempty"`
}

// Render the GolangConfig into its corresponding GolangResult.
func (c GolangConfig) Render() (Result, error) {
	l := log.WithField("src", "golang")
	l.Debug("starting render")

	result := NewGolangResult()

	for _, opt := range c.Core {
		switch opt {
		case "version":
			if !binExists("go") {
				cliWarnings.Add("go.core.version", "go executable not found")
				continue
			}
			stdout, stderr, err := runCommand("go", "version")
			if err != nil {
				errString := stderr.String()
				if errString == "" {
					errString = "<no output>"
				}
				l.Debugf("command error: %v", errString)
				cliWarnings.Add("go.core.version", "unable to determine version of go")
				continue
			}
			result.Version = toSlice(stdout.Bytes())[2]

		case "goroot":
			if !binExists("go") {
				cliWarnings.Add("go.core.goroot", "go executable not found")
				continue
			}
			stdout, stderr, err := runCommand("go", "env", "GOROOT")
			if err != nil {
				errString := stderr.String()
				if errString == "" {
					errString = "<no output>"
				}
				l.Debugf("command error: %v", errString)
				cliWarnings.Add("go.core.goroot", "unable to determine GOROOT")
				continue
			}
			result.Goroot = normalize(stdout.Bytes())

		case "gopath":
			if !binExists("go") {
				cliWarnings.Add("go.core.gopath", "go executable not found")
				continue
			}
			stdout, stderr, err := runCommand("go", "env", "GOPATH")
			if err != nil {
				errString := stderr.String()
				if errString == "" {
					errString = "<no output>"
				}
				l.Debugf("command error: %v", errString)
				cliWarnings.Add("go.core.gopath", "unable to determine GOPATH")
				continue
			}
			result.Gopath = normalize(stdout.Bytes())

		default:
			return result, fmt.Errorf("unsupported core golang option: %s", opt)
		}
	}

	return result, nil
}

// GolangResult contains the result data from rendering an "go" source.
type GolangResult struct {
	// Common
	ResultCommon `json:"-" yaml:"-"`

	// Core
	Version string `yaml:"version,omitempty" json:"version,omitempty"`
	Goroot  string `yaml:"goroot,omitempty" json:"goroot,omitempty"`
	Gopath  string `yaml:"gopath,omitempty" json:"gopath,omitempty"`
}

// NewGolangResult creates a new instance of an GolangResult.
func NewGolangResult() GolangResult {
	return GolangResult{
		ResultCommon: common,
	}
}

// IsEmpty checks whether the result contains any data.
func (r GolangResult) IsEmpty() bool {
	return r.Version == "" && r.Goroot == "" && r.Gopath == ""
}

// Markdown renders the GolangResult to markdown.
func (r GolangResult) Markdown() ([]byte, error) {
	log.WithField("src", "golang").Debug("rendering to markdown")

	if r.IsEmpty() {
		return []byte{}, nil
	}

	md := heredoc.Doc(`
		**Golang**{{ if .Version }}
		- _version_: {{ .Version }}{{ end }}{{ if .Goroot }}
		- _goroot_: {{ .Goroot }}{{ end }}{{ if .Gopath }}
		- _gopath_: {{ .Gopath }}{{ end }}
	`)
	t := template.Must(template.New("golang-md").Parse(md))

	buffer := bytes.Buffer{}
	if err := t.Execute(&buffer, r); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

// Plaintext renders the GolangResult to plaintext.
func (r GolangResult) Plaintext() ([]byte, error) {
	log.WithField("src", "golang").Debug("rendering to plaintext")

	if r.IsEmpty() {
		return []byte{}, nil
	}

	plaintext := heredoc.Doc(`
		Golang
		------{{ if .Version }}
		version:  {{ .Version }}{{ end }}{{ if .Goroot }}
		goroot:   {{ .Goroot }}{{ end }}{{ if .Gopath }}
		gopath:   {{ .Gopath }}{{ end }}
	`)

	t := template.Must(template.New("golang-txt").Parse(plaintext))

	buffer := bytes.Buffer{}
	if err := t.Execute(&buffer, r); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

// YAML renders the GolangResult to YAML.
func (r GolangResult) YAML() ([]byte, error) {
	log.WithField("src", "golang").Debug("rendering to YAML")

	if r.IsEmpty() {
		return []byte{}, nil
	}
	return yaml.Marshal(&r)
}

// JSON renders the GolangResult to JSON.
func (r GolangResult) JSON() ([]byte, error) {
	log.WithField("src", "golang").Debug("rendering to JSON")

	if r.IsEmpty() {
		return []byte{}, nil
	}
	return json.Marshal(&r)
}
