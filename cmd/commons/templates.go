package commons

// GetAppHelpTemplate provides general application help template
func GetAppHelpTemplate() string {
	return `NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}
USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}{{end}}{{end}}{{if .Description}}
DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}
AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}
COMMANDS:{{range .VisibleCategories}}{{if .Name}}

   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}} {{"\t"}}{{.Description}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

GLOBAL OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{.Copyright}}{{end}}
`
}

// GetCommandHelpTemplate provides general command help template
func GetCommandHelpTemplate() string {
	return `NAME:
   {{.HelpName}}{{if .Description}} - {{.Description}}{{end}}
   {{if .VisibleFlags}}
OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Usage}}
USAGE:
   {{.Usage}}{{end}}

`
}

// GetSubcommandHelpTemplate provides general sub command help template
func GetSubcommandHelpTemplate() string {
	return `NAME:
   {{.HelpName}} - {{if .Description}}{{.Description}}
   {{end}}
USAGE:
   {{.HelpName}} subcommand{{if .VisibleFlags}} [subcommand options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]
   {{end}}
SUBCOMMANDS:{{range .VisibleCommands}}{{if .Description}}
     {{join .Names ", "}}{{"\t"}}{{.Description}}{{end}}{{end}}

`
}
