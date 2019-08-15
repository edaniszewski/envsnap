package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// PythonConfig defines the configuration for the "python" source.
type PythonConfig struct {
	Core []string           `yaml:"core,omitempty"`
	Deps DependenciesConfig `yaml:"dependencies,omitempty"`
}

// DependenciesConfig defines the Python configuration for specifying
// package dependencies.
type DependenciesConfig struct {
	Packages []string `yaml:"packages,omitempty"`

	// todo: not yet implemented
	//From []string `yaml:"from"`
}

// Render the PythonConfig into its corresponding PythonResult.
func (c PythonConfig) Render() (Result, error) {
	l := log.WithField("src", "python")
	l.Debug("starting render")

	result := NewPythonResult()

	// Core Options
	for _, opt := range c.Core {
		switch opt {
		case "version":
			if !binExists("python") {
				cliWarnings.Add("python.core.version", "python executable not found")
				continue
			}

			stdout, stderr, err := runCommand("python", "--version")
			if err != nil {
				errString := stderr.String()
				if errString == "" {
					errString = "<no output>"
				}
				l.Debugf("command error: %v", errString)
				cliWarnings.Add("python.core.version", "unable to determine version of python")
				continue
			}
			result.Version = toSlice(stdout.Bytes())[1]

		case "py2":
			if !binExists("python2") {
				cliWarnings.Add("python.core.py2", "python2 executable not found")
				continue
			}
			stdout, stderr, err := runCommand("python2", "--version")
			if err != nil {
				errString := stderr.String()
				if errString == "" {
					errString = "<no output>"
				}
				l.Debugf("command error: %v", errString)
				cliWarnings.Add("python.core.py2", "unable to determine version of python2")
				continue
			}
			result.VersionPy2 = toSlice(stdout.Bytes())[1]

		case "py3":
			if !binExists("python3") {
				cliWarnings.Add("python.core.py3", "python3 executable not found")
				continue
			}
			stdout, stderr, err := runCommand("python3", "--version")
			if err != nil {
				errString := stderr.String()
				if errString == "" {
					errString = "<no output>"
				}
				l.Debugf("command error: %v", errString)
				cliWarnings.Add("python.core.py3", "unable to determine version of python3")
				continue
			}
			result.VersionPy3 = toSlice(stdout.Bytes())[1]

		default:
			return result, fmt.Errorf("unsupported option for python.core: %s", opt)
		}
	}

	// Dependencies Options
	if len(c.Deps.Packages) != 0 {
		if !binExists("pip") {
			cliWarnings.Add("python.deps.packages", "pip executable not found")
		} else {
			for _, dep := range c.Deps.Packages {

				stdout, stderr, err := runCommand("pip", "show", dep)
				if err != nil {
					errString := stderr.String()
					if errString == "" {
						errString = "<no output>"
					}
					l.WithField("dep", dep).Debugf("command error: %v", errString)
					cliWarnings.Add(
						"python.dependencies.packages",
						"python dependency not found: '%s'", dep,
					)
					result.Deps[dep] = ""
					continue
				}

				fields := strings.Split(string(stdout.String()), "\n")
				name := strings.Split(fields[0], " ")[1]
				version := strings.Split(fields[1], " ")[1]
				result.Deps[name] = version
			}
		}
	}

	// todo: not yet implemented
	//if len(c.Deps.From) != 0 {
	//	for _, source := range c.Deps.From {
	//		switch source {
	//		case "setup.py":
	//		case "requirements.txt":
	//		case "requirements.in:":
	//		default:
	//			return result, fmt.Errorf("unsupported option for python.dependencies.for: %s", source)
	//		}
	//	}
	//}

	return result, nil
}

// PythonResult contains the result data from rendering an "python" source.
type PythonResult struct {
	// Common
	ResultCommon `json:"-" yaml:"-"`

	// Core
	Version    string `yaml:"version,omitempty" json:"version,omitempty"`
	VersionPy2 string `yaml:"py2,omitempty" json:"py2,omitempty"`
	VersionPy3 string `yaml:"py3,omitempty" json:"py3,omitempty"`

	// Dependencies
	Deps map[string]string `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
}

// NewPythonResult creates a new instance of an PythonResult.
func NewPythonResult() PythonResult {
	return PythonResult{
		ResultCommon: common,
		Deps:         make(map[string]string),
	}
}

// IsEmpty checks whether the result contains any data.
func (r PythonResult) IsEmpty() bool {
	return r.Version == "" && r.VersionPy2 == "" && r.VersionPy3 == "" && len(r.Deps) == 0
}

// Markdown renders the PythonResult to markdown.
func (r PythonResult) Markdown() ([]byte, error) {
	log.WithField("src", "python").Debug("rendering to markdown")

	if r.IsEmpty() {
		return []byte{}, nil
	}

	md := heredoc.Doc(`
		**Python**{{ if .Version }}
		- _version_: {{ .Version }}{{ end }}{{ if .VersionPy2 }}
		- _py2_: {{ .VersionPy2 }}{{ end }}{{ if .VersionPy3 }}
		- _py3_: {{ .VersionPy3 }}{{ end }}{{ if .Deps }}
		- _dependencies_:
		  {{ .CodeFence }}
		  {{ range $key, $val := .Deps }}{{ $key }}=={{ $val }}
		  {{ end }}{{ .CodeFence }}{{ end }}
	`)
	t := template.Must(template.New("python-md").Parse(md))

	buffer := bytes.Buffer{}
	if err := t.Execute(&buffer, r); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

// Plaintext renders the PythonResult to plaintext.
func (r PythonResult) Plaintext() ([]byte, error) {
	log.WithField("src", "python").Debug("rendering to plaintext")

	if r.IsEmpty() {
		return []byte{}, nil
	}

	plaintext := heredoc.Doc(`
		Python
		------{{ if .Version }}
		version:  {{ .Version }}{{ end }}{{ if .VersionPy2 }}
		py2:      {{ .VersionPy2 }}{{ end }}{{ if .VersionPy3 }}
		py3:      {{ .VersionPy3 }}{{ end }}{{ if .Deps }}
		dependencies:
		{{ range $key, $val := .Deps }}- {{ $key }}=={{ $val }}
		{{ end }}{{ end -}}
	`)

	t := template.Must(template.New("python-txt").Parse(plaintext))

	buffer := bytes.Buffer{}
	if err := t.Execute(&buffer, r); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

// YAML renders the PythonResult to YAML.
func (r PythonResult) YAML() ([]byte, error) {
	log.WithField("src", "python").Debug("rendering to YAML")

	if r.IsEmpty() {
		return []byte{}, nil
	}
	return yaml.Marshal(&r)
}

// JSON renders the PythonResult to JSON.
func (r PythonResult) JSON() ([]byte, error) {
	log.WithField("src", "python").Debug("rendering to JSON")

	if r.IsEmpty() {
		return []byte{}, nil
	}
	return json.Marshal(&r)
}
