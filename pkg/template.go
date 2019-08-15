package pkg

import "github.com/MakeNowJust/heredoc"

// AppHelpTemplate is a custom template for the CLI help message.
var AppHelpTemplate = heredoc.Doc(`
	Usage: {{ .HelpName }} [global options] command [command options] [arguments ...]
	
	{{ .Usage }}
	
	Commands:{{range .VisibleCategories}}{{if .Name}}
	
	   {{.Name}}:{{range .VisibleCommands}}
		 {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
	   {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
	
	Options:
	   {{range $index, $option := .VisibleFlags}}{{if $index}}
	   {{end}}{{$option}}{{end}}{{end}}	
	
	Run '{{ .HelpName }} [command] --help' for more information about a command.
`)

// CommandHelpTemplate is a custom template for the CLI command help messages.
var CommandHelpTemplate = heredoc.Doc(`
	Usage: {{ .HelpName }}{{ if .VisibleFlags }} [options]{{ end }} {{ if .ArgsUsage }}{{ .ArgsUsage }}{{ else }}[arguments...]{{ end }}
	
	{{ .Usage }}{{ if .Description }}
	
	{{ .Description }}{{ end }}{{if .VisibleFlags}}
	
	Options:
	   {{range .VisibleFlags}}{{.}}
	   {{end}}{{end}}
`)

// CommandVersionTemplate is the template for the CLI version output.
var CommandVersionTemplate = heredoc.Doc(`
	envsnap:
	  version:      {{ .Version }}
	  build date:   {{ .BuildDate }}
	  commit:       {{ .Commit }}
	  tag:          {{ .Tag }}
	  go version:   {{ .GoVersion }}
	  go compiler:  {{ .Compiler }}
	  platform:     {{ .OS }}/{{ .Arch }}
`)

// EnvsnapInitTemplate is the template for the boilerplate envsnap config.
var EnvsnapInitTemplate = heredoc.Doc(`
	# envsnap configuration (yaml format)
	# use 'envsnap show' to generate an environment snapshot
	# for more details, see: https://www.github.com/edaniszewski/envsnap
	
	version: {{ .Version }}
	{{ if not .Terse }}
	# System configurations provide details about the user's system.{{ end }}
	system:
	  core:
	  - os
	  - arch
	{{ if .RenderPython -}}
	{{ if not .Terse }}
	# Python configurations provide details about the user's Python
	# installation and dependencies.{{ end }}
	python:
	  core:
	  - version
	  dependencies: []
	{{ end }}{{ if .RenderGolang -}}
	{{ if not .Terse }}
	# Golang configurations provide details about the user's Go installation.{{ end }}
	go:
	  core:
	  - version
	{{ end -}}
`)
