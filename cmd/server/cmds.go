package server

import (
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

// GetCommands returns slice of commands for this command category
func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "server",
			Aliases:     []string{"servers"},
			Hidden:      false,
			Description: "Group of commands related to servers",
			Subcommands: []cli.Command{
				{
					Name:    "add",
					Aliases: []string{"update"},
					Flags: []cli.Flag{
						commons.ServerFlag,
						cli.StringFlag{
							Name:  "api-url",
							Value: "",
							Usage: "Required. RabbitMQ api url, for example http://localhost:15672",
						},
						cli.StringFlag{
							Name:  "amqp-url",
							Value: "",
							Usage: "Required. RabbitMQ amqp url, for example amqp://localhost:5672",
						},
						cli.StringFlag{
							Name:  "username, u",
							Value: "",
							Usage: "Required. Username used during authentication to the RabbitMQ server",
						},
						cli.StringFlag{
							Name:  "password, p",
							Value: "",
							Usage: "Required. Password used during authentication to the RabbitMQ server",
						},
					},
					Description: "Add or update provided RabbitMQ server configuration",
					Action:      addServerCmd,
					Usage: "\n\t" +
						"rabbitr server add -s my-server-from-cfg -api-url http://localhost:15672 -amqp-url amqp://localhost:5672 -u user -p pass\t# will add new or update existing server to the configuration under the my-server-from-cfg name\n\t",
				},
				{
					Name:    "delete",
					Aliases: []string{"remove"},
					Flags: []cli.Flag{
						commons.ServerFlag,
					},
					Description: "Delete RabbitMQ server from the configuration",
					Action:      deleteServerConfigCmd,
					Usage: "\n\t" +
						"rabbitr server delete -s my-server-from-cfg \t# will delete server named my-server-from-cfg from the configuration\n\t",
				},
				{
					Name:        "list",
					Aliases:     []string{"show"},
					Description: "List all servers defined in configuration",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:   "show-passwords",
							Usage:  "Show passwords instead of redacted text",
							Hidden: true,
						},
					},
					Action: showConfigurationCmd,
					Usage: "\n\t" +
						"rabbitr servers list \t# will list servers defined in configuration\n\t",
				},
			},
		},
	}
}
