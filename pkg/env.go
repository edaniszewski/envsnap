package pkg

import (
	"bytes"
	"encoding/json"
	"os"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// EnvConfig defines the configuration for the "environment" source.
type EnvConfig struct {
	Variables []string `yaml:"variables,omitempty"`
}

// Render the EnvConfig into its corresponding EnvResult.
func (c EnvConfig) Render() (Result, error) {
	l := log.WithField("src", "env")
	l.Debug("starting render")

	result := NewEnvResult()

	for _, key := range c.Variables {
		val, _ := os.LookupEnv(key)
		l.WithFields(log.Fields{
			"key": key,
			"val": val,
		}).Debug("env lookup")
		result.Env[key] = val
	}
	return result, nil
}

// EnvResult contains the result data from rendering an "environment" source.
type EnvResult struct {
	// Common
	ResultCommon `json:"-" yaml:"-"`

	// Env
	Env map[string]string `yaml:",omitempty,inline" json:"env,omitempty"`
}

// NewEnvResult creates a new instance of an EnvResult.
func NewEnvResult() EnvResult {
	return EnvResult{
		ResultCommon: common,
		Env:          make(map[string]string),
	}
}

// IsEmpty checks whether the result contains any data.
func (r EnvResult) IsEmpty() bool {
	return len(r.Env) == 0
}

// Markdown renders the EnvResult to markdown.
func (r EnvResult) Markdown() ([]byte, error) {
	log.WithField("src", "env").Debug("rendering to markdown")

	if r.IsEmpty() {
		return []byte{}, nil
	}

	md := heredoc.Doc(`
		**Environment**
		{{ .CodeFence }}
		{{ range $k, $v := .Env }}{{ $k }}={{ $v }}
		{{ end }}{{ .CodeFence }}
	`)
	t := template.Must(template.New("env-md").Parse(md))

	buffer := bytes.Buffer{}
	if err := t.Execute(&buffer, r); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

// Plaintext renders the EnvResult to plaintext.
func (r EnvResult) Plaintext() ([]byte, error) {
	log.WithField("src", "env").Debug("rendering to plaintext")

	if r.IsEmpty() {
		return []byte{}, nil
	}

	plaintext := heredoc.Doc(`
		Environment
		-----------
		{{ range $k, $v := .Env }}{{ $k }}={{ $v }}
		{{ end -}}
	`)

	t := template.Must(template.New("env-txt").Parse(plaintext))

	buffer := bytes.Buffer{}
	if err := t.Execute(&buffer, r); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

// YAML renders the EnvResult to YAML.
func (r EnvResult) YAML() ([]byte, error) {
	log.WithField("src", "env").Debug("rendering to YAML")

	if r.IsEmpty() {
		return []byte{}, nil
	}
	return yaml.Marshal(&r)
}

// JSON renders the EnvResult to JSON.
func (r EnvResult) JSON() ([]byte, error) {
	log.WithField("src", "env").Debug("rendering to JSON")

	if r.IsEmpty() {
		return []byte{}, nil
	}
	return json.Marshal(&r)
}
