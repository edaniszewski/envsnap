package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var common ResultCommon

func init() {
	common = ResultCommon{
		CodeQuote: "`",
		CodeFence: "```",
	}
}

// ResultCommon is a struct which defines common components to results which
// instances of a Result may embed to access the fields here for output formatting,
// particularly for markdown.
type ResultCommon struct {
	// TODO move somewhere outside of interface
	CodeQuote string
	CodeFence string
}

// Result defines an interface for rendered configurations which allows
// them to be output in a number of supported formats.
type Result interface {
	IsEmpty() bool

	Markdown() ([]byte, error)
	Plaintext() ([]byte, error)
	YAML() ([]byte, error)
	JSON() ([]byte, error)
}

// EnvsnapResult defines an interface which all versions of envsnap results
// should implement.
type EnvsnapResult interface {
	Results() []Result
	String(format string) (string, error)
	Write(file, format string) error
	Print(format string) error
}

// V1EnvsnapResult contains the results for all sources specified by version 1
// of the envsnap configuration, as defined in V1EnvsnapConfig.
type V1EnvsnapResult struct {
	Environment Result `json:"environment,omitempty" yaml:"environment,omitempty"`
	Exec        Result `json:"exec,omitempty" yaml:"exec,omitempty"`
	Golang      Result `json:"golang,omitempty" yaml:"golang,omitempty"`
	Python      Result `json:"python,omitempty" yaml:"python,omitempty"`
	System      Result `json:"system,omitempty" yaml:"system,omitempty"`

	out io.Writer
}

// NewV1EnvsnapResult creates a new instance of the V1EnvsnapResult struct, setting
// the default `out` to stdout.
func NewV1EnvsnapResult() V1EnvsnapResult {
	return V1EnvsnapResult{
		out: os.Stdout,
	}
}

// Results returns all of the component source results in the order in which
// they should be rendered.
func (r *V1EnvsnapResult) Results() []Result {
	return []Result{
		r.System,
		r.Environment,
		r.Exec,
		r.Python,
		r.Golang,
	}
}

// String renders the result into a string based on the given format option.
// If the provided format is not supported, this returns an error.
func (r *V1EnvsnapResult) String(format string) (string, error) {
	switch format {
	case "markdown", "md":
		var parts []string
		for _, res := range r.Results() {
			if res == nil || res.IsEmpty() {
				continue
			}
			data, err := res.Markdown()
			if err != nil {
				return "", err
			}
			parts = append(parts, string(data))
		}
		return "#### Environment\n\n" + strings.Join(parts, "\n"), nil

	case "plaintext", "txt":
		var parts []string
		for _, res := range r.Results() {
			if res == nil || res.IsEmpty() {
				continue
			}
			data, err := res.Plaintext()
			if err != nil {
				return "", err
			}
			parts = append(parts, string(data))
		}
		return strings.Join(parts, "\n"), nil

	case "yaml":
		data, err := yaml.Marshal(r)
		if err != nil {
			return "", err
		}
		return string(data), nil

	case "json":
		data, err := json.Marshal(r)
		if err != nil {
			return "", err
		}
		return string(data), nil

	default:
		return "", ErrUnsupportedFormat
	}
}

// Write renders the result into a string based on the provided format and
// writes that string to the specified file.
func (r *V1EnvsnapResult) Write(file, format string) error {
	out, err := r.String(format)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(out), 0644)
}

// Print renders the result into a string based on the provided format and
// writes that string to stdout.
func (r *V1EnvsnapResult) Print(format string) error {
	out, err := r.String(format)
	if err != nil {
		return err
	}
	fmt.Fprintln(r.out, out)
	return nil
}
