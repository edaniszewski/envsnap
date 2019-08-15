package pkg

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ExecConfig defines the configuration for the "exec" source.
type ExecConfig struct {
	Run []string `yaml:"run,omitempty"`
}

// Render the ExecConfig into its corresponding ExecResult.
func (c ExecConfig) Render() (Result, error) {
	l := log.WithField("src", "exec")
	l.Debug("starting render")

	result := NewExecResult()

	for _, cmdStr := range c.Run {
		args := strings.Split(cmdStr, " ")
		l.WithField("cmd", cmdStr).Debug("running command")
		stdout, stderr, err := runCommand(args[0], args[1:]...)
		if err != nil {
			errString := stderr.String()
			if errString == "" {
				errString = "<no output>"
			}
			l.Debugf("command error: %v", errString)
			cliWarnings.Add(
				"exec.run",
				"error while running command: '%s'", cmdStr,
			)
			result.Exec[cmdStr] = ""
			continue
		}
		result.Exec[cmdStr] = stdout.String()
	}

	return result, nil
}

// ExecResult contains the result data from rendering an "exec" source.
type ExecResult struct {
	// Common
	ResultCommon `json:"-" yaml:"-"`

	// Exec
	Exec map[string]string `yaml:"exec,omitempty" json:"exec,omitempty"`
}

// NewExecResult creates a new instance of an ExecResult.
func NewExecResult() ExecResult {
	return ExecResult{
		ResultCommon: common,
		Exec:         make(map[string]string),
	}
}

// IsEmpty checks whether the result contains any data.
func (r ExecResult) IsEmpty() bool {
	return len(r.Exec) == 0
}

// Markdown renders the ExecResult to markdown.
func (r ExecResult) Markdown() ([]byte, error) {
	log.WithField("src", "exec").Debug("rendering to markdown")

	if r.IsEmpty() {
		return []byte{}, nil
	}

	md := heredoc.Doc(`
		**Exec**{{ if .Exec }}{{ $root := . }}
		{{ range $key, $val := .Exec }}- {{ $root.CodeQuote }}{{ $key }}{{ $root.CodeQuote }}
		  {{ $root.CodeFence }}
		  {{ $val }}
		  {{ $root.CodeFence }}
		{{ end }}{{ end -}}
	`)
	t := template.Must(template.New("exec-md").Parse(md))

	buffer := bytes.Buffer{}
	if err := t.Execute(&buffer, r); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

// Plaintext renders the ExecResult to plaintext.
func (r ExecResult) Plaintext() ([]byte, error) {
	log.WithField("src", "exec").Debug("rendering to plaintext")

	if r.IsEmpty() {
		return []byte{}, nil
	}

	plaintext := heredoc.Doc(`
		Exec
		----{{ if .Exec }}
		{{ range $key, $val := .Exec }}$ {{ $key }}
		  {{ $val }}
		{{ end }}{{ end -}}
	`)

	t := template.Must(template.New("exec-txt").Parse(plaintext))

	buffer := bytes.Buffer{}
	if err := t.Execute(&buffer, r); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

// YAML renders the ExecResult to YAML.
func (r ExecResult) YAML() ([]byte, error) {
	log.WithField("src", "exec").Debug("rendering to YAML")

	if r.IsEmpty() {
		return []byte{}, nil
	}
	return yaml.Marshal(&r)
}

// JSON renders the ExecResult to JSON.
func (r ExecResult) JSON() ([]byte, error) {
	log.WithField("src", "exec").Debug("rendering to JSON")

	if r.IsEmpty() {
		return []byte{}, nil
	}
	return json.Marshal(&r)
}
