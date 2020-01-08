package policy

import (
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "policy",
			Aliases: []string{"policies"},
			Hidden:  false,
			Subcommands: []cli.Command{
				{
					Name: "list",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.VHostFlag,
						commons.GetFilterFlag(commons.AllSubjects, "policy.Name/Vhost/Pattern/ApplyTo/Priority"),
					},
					Description: "Lists policies defined on RabbitMQ server",
					Action:      listPoliciesCmd,
					Usage: "\n\t" +
						"rabbitr policies list -s my-server-from-cfg\t# will list all policies defined in my-server-from-cfg\n\t" +
						"rabbitr policies list -s my-server-from-cfg -f 'policy.Name=~\"Federation\"'\t# will list policies defined on my-server-from-cfg which's name matches the 'Federation' string\n\t" +
						"rabbitr policies list -s my-server-from-cfg -f 'policy.Priority>8'\t# will list policies defined on my-server-from-cfg with priorities greater than 8\n\t",
				},
			},
		},
	}
}
