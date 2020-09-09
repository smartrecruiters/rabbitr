package exchange

import (
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "exchange",
			Aliases: []string{"exchanges"},
			Hidden:  false,
			Subcommands: []cli.Command{
				{
					Name: "list",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.VHostFlag,
						commons.GetFilterFlag(commons.AllSubjects, commons.ExchangeFilterFields),
					},
					Description: "Lists exchanges on given RabbitMQ server",
					Action:      listExchangesCmd,
					Usage: "\n\t" +
						"rabbitr exchanges list -s my-server-from-cfg\t# will list all exchanges from my-server-from-cfg\n\t" +
						"rabbitr exchanges list -s my-server-from-cfg -f 'exchange.Name=~\"my-exchange\"'\t# will list exchanges from my-server-from-cfg which's name matches the 'my-exchange' string\n\t",
				},
				{
					Name: "delete",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.DryRunFlag,
						commons.VHostFlag,
						commons.GetFilterFlag(commons.NoneOfTheSubjects, commons.ExchangeFilterFields),
					},
					Description: "Deletes exchange on given RabbitMQ server",
					Action:      deleteCmd,
					Usage: "\n\t" +
						"rabbitr exchange delete -s my-server-from-cfg -f \"1==1\"\t# will delete all exchanges from my-server-from-cfg\n\t" +
						"rabbitr exchange delete -s my-server-from-cfg -f 'exchange.Name=~\"my-exchange\"'\t# will delete exchanges from my-server-from-cfg which's name matches the 'my-exchange' string\n\t",
				},
			},
		},
	}
}
