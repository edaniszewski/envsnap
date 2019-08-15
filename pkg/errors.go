package pkg

import "errors"

// Errors used throughout envsnap.
var (
	ErrConfigExists         = errors.New(".envsnap file already exists")
	ErrNoConfig             = errors.New(".envsnap file not found")
	ErrUnsupportedLang      = errors.New("unsupported language passed to the --lang flag")
	ErrIncompleteRender     = errors.New("envsnap failed to render some configured options (run with --debug for more detail)")
	ErrUnsupportedFormat    = errors.New("unsupported format string provided")
	ErrNoConfigVersion      = errors.New("no version specified in config")
	ErrInvalidConfigVersion = errors.New("invalid config version specified")
	ErrInvalidGithubURL     = errors.New("invalid github url: must be in the format 'github.com/<user>/<repo>'")
)
