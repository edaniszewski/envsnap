package pkg

import (
	"os"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// InitOptions are the options used to render the .envsnap config boilerplate
// via `envsnap init`.
type InitOptions struct {
	Version int
	Terse   bool

	RenderPython bool
	RenderGolang bool
}

// NewApp creates a new instance of the envsnap CLI application.
func NewApp() *cli.App {

	// Set custom template for command help
	cli.CommandHelpTemplate = CommandHelpTemplate

	// Set a custom version printer
	cli.VersionPrinter = func(c *cli.Context) {
		t := template.Must(template.New("version").Parse(CommandVersionTemplate))
		if err := t.Execute(os.Stdout, newVersionInfo()); err != nil {
			log.Fatal(err)
		}
	}

	app := cli.NewApp()
	app.Name = "envsnap"
	app.Usage = "Generate project-defined snapshots of runtime environments"
	app.Version = Version
	app.CustomAppHelpTemplate = AppHelpTemplate
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "run envsnap with debug logging",
		},
	}
	app.Before = func(context *cli.Context) error {
		if context.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		} else {
			// If debug logging isn't enabled, set it to panic level to effectively
			// disable logging.
			log.SetLevel(log.PanicLevel)
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Initialize a boilerplate .envsnap config",
			Description: heredoc.Doc(`
				Initialize a boilerplate .envsnap config.

				The .envsnap file is a YAML-formatted file used by the envsnap CLI tool to
				generate a snapshot of the specified environment.

				The boilerplate config will contain section comments by default. This
				can be disabled with the '--terse' flag.

				To add language-specific sections to the generated config, use the '--lang'
				flag, passing to it the language(s), or language shorthand(s), you wish to
				include.

				Currently, the supported languages (shorthands in parentheses) are:
				  • python (py)
				  • golang (go)
				`,
			),
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "terse, t",
					Usage: "initialize the config without extra comments",
				},
				cli.StringSliceFlag{
					Name:  "lang, l",
					Usage: "initialize the config with basic sections for the specified languages",
				},
			},
			Action: commandInit,
		},
		{
			Name:  "render",
			Usage: "Render the environment as specified by the config",
			Description: heredoc.Doc(`
				The rendered environment is output to console by default. The '--file' flag can
				be used to write the output to file.

				The output format can be set with the '--output' flag. By default, it will render
				the results in markdown format. The allowable output formats are:
				  • md		Markdown output  (.md)
				  • txt		Plaintext output (.txt)
				  • yaml	YAML output      (.yaml)
				  • json	JSON output      (.json)
				`,
			),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Value: "md",
					Usage: "specify the output format",
				},
				cli.StringFlag{
					Name:  "file, f",
					Usage: "write the output to file",
				},
				cli.BoolFlag{
					Name:  "quiet, q",
					Usage: "ignore any warnings generated during render",
				},
			},
			Action: commandRender,
		},
	}

	return app
}
