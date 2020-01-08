package connection

import (
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "connection",
			Aliases: []string{"connections"},
			Hidden:  false,
			Subcommands: []cli.Command{
				{
					Name: "list",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.VHostFlag,
						commons.GetFilterFlag(commons.AllSubjects, "connection.ID/Name/Vhost"),
					},
					Description: "Lists connections made to RabbitMQ server",
					Action:      listConnectionsCmd,
					Usage: "\n\t" +
						"rabbitr connections list -s my-server-from-cfg\t# will list all connections made to my-server-from-cfg\n\t" +
						"rabbitr connections list -s my-server-from-cfg -f 'connection.Name=~\"Federation\"'\t# will list connections made to my-server-from-cfg which's name matches the 'Federation' string\n\t" +
						"rabbitr connections list -s my-server-from-cfg -f 'connection.ID=~\"10.30\"'\t# will list connections made to my-server-from-cfg which's id matches the '10.30' string\n\t",
				},
				{
					Name: "close",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.DryRunFlag,
						commons.VHostFlag,
						commons.GetFilterFlag(commons.NoneOfTheSubjects, "connection.ID/Name/Vhost"),
					},
					Description: "Closes connections made to RabbitMQ server",
					Action:      closeConnectionsCmd,
					Usage: "\n\t" +
						"rabbitr connections close -s my-server-from-cfg -dry-run=true\t# will list all connections made to my-server-from-cfg without closing them\n\t" +
						"rabbitr connections close -s my-server-from-cfg -f 'connection.Name=~\"Federation\"'\t# will close connections made to my-server-from-cfg which's name matches the 'Federation' string\n\t" +
						"rabbitr connections close -s my-server-from-cfg -f 'connection.ID=~\"10.27.160.18 ->\"'\t# will close connections made to my-server-from-cfg which's id matches the '10.27.160.18 ->' string\n\t",
				},
			},
		},
	}
}
