package main

import (
	"fmt"
	"os"

	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/smartrecruiters/rabbitr/cmd/connection"
	"github.com/smartrecruiters/rabbitr/cmd/exchange"
	"github.com/smartrecruiters/rabbitr/cmd/policy"
	"github.com/smartrecruiters/rabbitr/cmd/queue"
	"github.com/smartrecruiters/rabbitr/cmd/server"
	"github.com/smartrecruiters/rabbitr/cmd/shovel"
	"github.com/urfave/cli"
)

const (
	applicationName        = "rabbitr"
	applicationDescription = "CLI application for easier management of RabbitMQ related tasks"
)

var (
	version string
	commit  string
	date    string
)

func versionString() string {
	return fmt.Sprintf("%s, commit %s, built at %s", version, commit, date)
}

func main() {
	app := cli.NewApp()
	app.Name = applicationName
	app.Usage = applicationDescription
	app.Version = versionString()
	app.Commands = connection.GetCommands()
	app.Commands = append(app.Commands, exchange.GetCommands()...)
	app.Commands = append(app.Commands, policy.GetCommands()...)
	app.Commands = append(app.Commands, queue.GetCommands()...)
	app.Commands = append(app.Commands, server.GetCommands()...)
	app.Commands = append(app.Commands, shovel.GetCommands()...)

	cli.AppHelpTemplate = commons.GetAppHelpTemplate()
	cli.CommandHelpTemplate = commons.GetCommandHelpTemplate()
	cli.SubcommandHelpTemplate = commons.GetSubcommandHelpTemplate()

	err := app.Run(os.Args)
	commons.AbortIfError(err)
}
