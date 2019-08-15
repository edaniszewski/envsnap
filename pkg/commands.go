package pkg

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/urfave/cli"
)

// commandInit is the function executed for the CLI's "init" command.
func commandInit(c *cli.Context) error {
	// If no path is provided, assume the current working directory.
	path := c.Args().Get(0)
	if path == "" {
		path = "."
	}

	path, err := filepath.Abs(filepath.Join(path, configFile))
	if err != nil {
		return err
	}

	// If the config file already exists, there is nothing to init.
	if _, err := os.Stat(path); err == nil {
		return ErrConfigExists
	}

	opts := InitOptions{
		Version: configV1,
		Terse:   c.Bool("terse"),
	}
	for _, lang := range c.StringSlice("lang") {
		switch lang {
		case "python", "py":
			opts.RenderPython = true
		case "golang", "go":
			opts.RenderGolang = true
		default:
			return ErrUnsupportedLang
		}
	}

	tmpl, err := template.New("init").Parse(EnvsnapInitTemplate)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, opts); err != nil {
		return err
	}

	return ioutil.WriteFile(path, out.Bytes(), 0644)
}

// commandRender is the function executed for the CLI's "render" command.
func commandRender(c *cli.Context) error {
	// If no path is provided, assume current working directory.
	path := c.Args().Get(0)

	// Get command flags.
	flagOutput := c.String("output")
	flagFile := c.String("file")

	cfg, err := LoadConfig(path)
	if err != nil {
		return err
	}

	res, err := cfg.Render()
	if err != nil {
		return err
	}

	if flagFile != "" {
		if err := res.Write(flagFile, flagOutput); err != nil {
			return err
		}
	} else {
		if err := res.Print(flagOutput); err != nil {
			return err
		}
	}

	// Check for any warnings and print them out.
	if !c.Bool("quiet") {
		if cliWarnings.HasWarnings() {
			cliWarnings.Print(os.Stderr)
			return ErrIncompleteRender
		}
	}
	return nil
}
