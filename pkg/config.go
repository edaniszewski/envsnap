package pkg

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
)

// The name of the config file that envsnap reads from.
const configFile = ".envsnap"

var (
	// Version 1 of the envsnap configuration file scheme.
	configV1 = 1
)

// RenderConfig defines an interface for configuration sections for envsnap
// which can be rendered into results.
type RenderConfig interface {
	Render() (Result, error)
}

// VersionedConfig is an intermediary struct which is used to load the
// version information from configuration YAML. This allows envsnap to
// determine the version of the configuration, and thus the correct struct
// to load the configuration into.
type VersionedConfig struct {
	Version *int `yaml:"version"`
}

// EnvsnapConfig defines an interface which all versions of the envsnap
// configuration should implement.
type EnvsnapConfig interface {
	All() []RenderConfig
	Render() (EnvsnapResult, error)
}

// LoadConfig loads the configuration for envsnap to render.
//
// If no path is specified, it assumes the configuration is in the current
// working directory. The path may also be a reference to a GitHub repository
// containing the configuration, in the form of "github.com/<owner>/<repo>[@<ref>]"
// where the <ref> may be a branch name, tag, or commit.
func LoadConfig(path string) (EnvsnapConfig, error) {
	var (
		err  error
		data []byte
	)

	// If no path is specified, assume the current working directory.
	if path == "" {
		path, err = filepath.Abs(filepath.Join(".", configFile))
		if err != nil {
			return nil, err
		}
	}

	// If the path starts with "github.com/", assume that it is referencing
	// a remote GitHub repo. This will cause the config to be loaded from
	// .envsnap file in that repo, instead of locally.
	if strings.HasPrefix(path, "github.com/") {
		parts := strings.Split(path, "@")
		var ref string
		url := parts[0]
		if len(parts) > 1 {
			ref = parts[1]
		}

		parts = strings.Split(url, "/")
		if len(parts) != 3 {
			return nil, ErrInvalidGithubURL
		}
		user := parts[1]
		repo := parts[2]

		ghc := github.NewClient(nil)
		content, _, _, err := ghc.Repositories.GetContents(
			context.Background(),
			user,
			repo,
			configFile,
			&github.RepositoryContentGetOptions{
				Ref: ref,
			},
		)
		if err != nil {
			return nil, err
		}
		ctnt, err := content.GetContent()
		if err != nil {
			return nil, err
		}
		data = []byte(ctnt)
	} else {
		// Load from file
		// Attempt to load the config from file. First, check that
		// the specified path even exists.
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, ErrNoConfig
		}

		// The file exists. Read from it and figure out which version
		// of the configuration it is.
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		data = contents
	}

	// ----
	v := &VersionedConfig{}
	if err := yaml.Unmarshal(data, v); err != nil {
		return nil, err
	}
	if v.Version == nil {
		return nil, ErrNoConfigVersion
	}

	switch *v.Version {
	case configV1:
		cfg := &V1EnvsnapConfig{}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
		return cfg, nil

	default:
		return nil, ErrInvalidConfigVersion
	}
}

// V1EnvsnapConfig contains all the data for the environment snapshot.
type V1EnvsnapConfig struct {
	Environment EnvConfig    `yaml:"environment,omitempty"`
	Exec        ExecConfig   `yaml:"exec,omitempty"`
	Golang      GolangConfig `yaml:"go,omitempty"`
	Python      PythonConfig `yaml:"python,omitempty"`
	System      SystemConfig `yaml:"system,omitempty"`
}

// All returns all of the configuration components for the v1 envsnap config.
func (c V1EnvsnapConfig) All() []RenderConfig {
	return []RenderConfig{
		c.System,
		c.Environment,
		c.Exec,
		c.Python,
		c.Golang,
	}
}

// Render each configured source into its corresponding v1 result.
func (c V1EnvsnapConfig) Render() (EnvsnapResult, error) {
	var err error
	v1 := NewV1EnvsnapResult()

	v1.System, err = c.System.Render()
	if err != nil {
		return nil, err
	}

	v1.Environment, err = c.Environment.Render()
	if err != nil {
		return nil, err
	}

	v1.Exec, err = c.Exec.Render()
	if err != nil {
		return nil, err
	}

	v1.Python, err = c.Python.Render()
	if err != nil {
		return nil, err
	}

	v1.Golang, err = c.Golang.Render()
	if err != nil {
		return nil, err
	}

	return &v1, nil
}
