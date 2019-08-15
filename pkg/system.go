package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// SysInfo contains information about the system.
type SysInfo struct {
	OS            string
	Kernel        string
	KernelVersion string
	Arch          string
	Processor     string
}

// SystemConfig defines the configuration for the "system" source.
type SystemConfig struct {
	Core []string `yaml:"core,omitempty"`
}

// Render the SystemConfig into its corresponding SystemResult.
func (c SystemConfig) Render() (Result, error) {
	l := log.WithField("src", "system")
	l.Debug("starting render")

	result := NewSystemResult()

	info, err := LoadSystemInfo()
	if err != nil {
		l.WithField("err", err).Debug("error collecting system info")
		cliWarnings.Add("system.core", "error collecting system info")
	}

	for _, opt := range c.Core {
		switch opt {
		case "os":
			result.OS = info.OS
		case "arch":
			result.Arch = info.Arch
		case "cpus":
			result.CPUs = runtime.NumCPU()
		case "kernel":
			result.Kernel = info.Kernel
		case "kernel_version", "kernel-version":
			result.KernelVersion = info.KernelVersion
		case "processor":
			result.Processor = info.Processor
		default:
			l.WithField("opt", opt).Debug("unsupported core system option")
			return result, fmt.Errorf("unsupported core system option: %s", opt)
		}
	}
	return result, nil
}

// SystemResult contains the result data from rendering an "system" source.
type SystemResult struct {
	// Common
	ResultCommon `json:"-" yaml:"-"`

	// Core
	OS            string `yaml:"os,omitempty" json:"os,omitempty"`
	Arch          string `yaml:"arch,omitempty" json:"arch,omitempty"`
	CPUs          int    `yaml:"cpus,omitempty" json:"cpus,omitempty"`
	Kernel        string `yaml:"kernel,omitempty" json:"kernel,omitempty"`
	KernelVersion string `yaml:"kernel_version,omitempty" json:"kernel_version,omitempty"`
	Processor     string `yaml:"processor,omitempty" json:"processor,omitempty"`
}

// NewSystemResult creates a new instance of an SystemResult.
func NewSystemResult() SystemResult {
	return SystemResult{
		ResultCommon: common,
	}
}

// IsEmpty checks whether the result contains any data.
func (r SystemResult) IsEmpty() bool {
	return r.OS == "" && r.Arch == "" && r.CPUs == 0 && r.KernelVersion == "" && r.Kernel == "" && r.Processor == ""
}

// Markdown renders the SystemResult to markdown.
func (r SystemResult) Markdown() ([]byte, error) {
	log.WithField("src", "system").Debug("rendering to markdown")

	if r.IsEmpty() {
		return []byte{}, nil
	}

	md := heredoc.Doc(`
		**System**{{ if .OS }}
		- _os_: {{ .OS }}{{ end }}{{ if .Arch }}
		- _arch_: {{ .Arch }}{{ end }}{{ if .CPUs }}
		- _cpus_: {{ .CPUs }}{{ end }}{{ if .Kernel }}
		- _kernel_: {{ .Kernel }}{{ end }}{{ if .KernelVersion }}
		- _kernel version_: {{ .KernelVersion }}{{ end }}{{ if .Processor }}
		- _processor_: {{ .Processor }}{{ end }}
	`)
	t := template.Must(template.New("system-md").Parse(md))

	buffer := bytes.Buffer{}
	if err := t.Execute(&buffer, r); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

// Plaintext renders the SystemResult to plaintext.
func (r SystemResult) Plaintext() ([]byte, error) {
	log.WithField("src", "system").Debug("rendering to plaintext")

	if r.IsEmpty() {
		return []byte{}, nil
	}

	plaintext := heredoc.Doc(`
		System
		------{{ if .OS }}
		os:             {{ .OS }}{{ end }}{{ if .Arch }}
		arch:           {{ .Arch }}{{ end }}{{ if .CPUs }}
		cpus:           {{ .CPUs }}{{ end }}{{ if .Kernel }}
		kernel:         {{ .Kernel }}{{ end }}{{ if .KernelVersion }}
		kernel version: {{ .KernelVersion }}{{ end }}{{ if .Processor }}
		processor:      {{ .Processor }}{{ end }}
	`)

	t := template.Must(template.New("system-txt").Parse(plaintext))

	buffer := bytes.Buffer{}
	if err := t.Execute(&buffer, r); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

// YAML renders the SystemResult to YAML.
func (r SystemResult) YAML() ([]byte, error) {
	log.WithField("src", "system").Debug("rendering to YAML")

	if r.IsEmpty() {
		return []byte{}, nil
	}
	return yaml.Marshal(&r)
}

// JSON renders the SystemResult to JSON.
func (r SystemResult) JSON() ([]byte, error) {
	log.WithField("src", "system").Debug("rendering to JSON")

	if r.IsEmpty() {
		return []byte{}, nil
	}
	return json.Marshal(&r)
}
