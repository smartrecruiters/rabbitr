package shovel

import (
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "shovel",
			Aliases: []string{"shovels"},
			Hidden:  false,
			Subcommands: []cli.Command{
				{
					Name: "list",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.VHostFlag,
						commons.GetFilterFlag(commons.AllSubjects, commons.ShovelFilterFields),
					},
					Description: "Lists shovels on given RabbitMQ server",
					Action:      listShovelsCmd,
					Usage: "\n\t" +
						"rabbitr shovels list -s my-server-from-cfg\t# will list all shovels from my-server-from-cfg\n\t" +
						"rabbitr shovels list -s my-server-from-cfg -f 'shovel.Name=~\"my-shovel\"'\t# will list shovels from my-server-from-cfg which's name matches the 'my-shovel' string\n\t",
				},
			},
		},
	}
}
