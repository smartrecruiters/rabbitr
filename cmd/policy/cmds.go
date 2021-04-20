package policy

import (
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

// GetCommands returns slice of commands for this command category
func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "policy",
			Aliases:     []string{"policies"},
			Hidden:      false,
			Description: "Group of commands related to policies",
			Subcommands: []cli.Command{
				{
					Name: "list",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.VHostFlag,
						commons.GetFilterFlag(commons.AllSubjects, commons.PolicyFilterFields),
					},
					Description: "Lists policies defined on RabbitMQ server",
					Action:      listPoliciesCmd,
					Usage: "\n\t" +
						"rabbitr policies list -s my-server-from-cfg\t# will list all policies defined in my-server-from-cfg\n\t" +
						"rabbitr policies list -s my-server-from-cfg -f 'policy.Name=~\"Federation\"'\t# will list policies defined on my-server-from-cfg which's name matches the 'Federation' string\n\t" +
						"rabbitr policies list -s my-server-from-cfg -f 'policy.Priority>8'\t# will list policies defined on my-server-from-cfg with priorities greater than 8\n\t",
				},
				{
					Name: "delete",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.DryRunFlag,
						commons.VHostFlag,
						commons.GetFilterFlag(commons.NoneOfTheSubjects, commons.PolicyFilterFields),
					},
					Description: "Deletes policy on given RabbitMQ server",
					Action:      deleteCmd,
					Usage: "\n\t" +
						"rabbitr policy delete -s my-server-from-cfg -f \"1==1\"\t# will delete all policies from my-server-from-cfg\n\t" +
						"rabbitr policy delete -s my-server-from-cfg -f 'policy.Name=~\"my-policy\"'\t# will delete policy from my-server-from-cfg which name matches the 'my-policy' string\n\t",
				},
			},
		},
	}
}
