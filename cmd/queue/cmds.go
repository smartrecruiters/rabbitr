package queue

import (
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

// GetCommands returns slice of commands for this command category
func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "queue",
			Aliases:     []string{"queues"},
			Hidden:      false,
			Description: "Group of commands related to queues",
			Subcommands: []cli.Command{
				{
					Name: "list",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.VHostFlag,
						commons.GetFilterFlag(commons.AllSubjects, commons.QueueFilterFields),
					},
					Description: "Lists queues on given RabbitMQ server",
					Action:      listQueuesCmd,
					Usage: "\n\t" +
						"rabbitr queues list -s my-server-from-cfg\t# will list all queues from my-server-from-cfg\n\t" +
						"rabbitr queues list -s my-server-from-cfg -f 'queue.Name=~\"my-queue\"'\t# will list queues from my-server-from-cfg which's name matches the 'my-queue' string\n\t" +
						"rabbitr queues list -s my-server-from-cfg -f 'queue.Messages>200'\t# will list queues from to my-server-from-cfg with more than 200 messages\n\t" +
						"rabbitr queues list -s my-server-from-cfg -f 'queue.Consumers==0'\t# will list queues from to my-server-from-cfg that have 0 consumers\n\t" +
						"rabbitr queues list -s my-server-from-cfg -f 'getMapValueByKey(queue.Arguments,\"x-max-priority\")==10'\t# will list queues from to my-server-from-cfg that have x-max-priority=10 defined in the queue arguments\n\t" +
						"rabbitr queues list -s my-server-from-cfg -f 'getMapValueByKey(queue.Arguments,\"x-queue-owner\")==\"queue@owner.com\"'\t# will list queues from to my-server-from-cfg that have x-queue-owner='queue@owner.com' defined in the queue arguments\n\t",
				},
				{
					Name: "delete",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.DryRunFlag,
						commons.VHostFlag,
						commons.GetFilterFlag(commons.NoneOfTheSubjects, commons.QueueFilterFields),
					},
					Description: "Deletes queues on given RabbitMQ server",
					Action:      deleteCmd,
					Usage: "\n\t" +
						"rabbitr queues delete -s my-server-from-cfg -f \"1==1\"\t# will delete all queues from my-server-from-cfg\n\t" +
						"rabbitr queues delete -s my-server-from-cfg -f 'queue.Name=~\"my-queue\"'\t# will delete queues from my-server-from-cfg which's name matches the 'my-queue' string\n\t" +
						"rabbitr queues delete -s my-server-from-cfg -f 'queue.Messages>200'\t# will delete queues from to my-server-from-cfg with more than 200 messages\n\t" +
						"rabbitr queues delete -s my-server-from-cfg -f 'queue.Consumers==0'\t# will delete queues from to my-server-from-cfg that have 0 consumers\n\t",
				},
				{
					Name: "purge",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.DryRunFlag,
						commons.VHostFlag,
						commons.GetFilterFlag(commons.NoneOfTheSubjects, commons.QueueFilterFields),
					},
					Description: "Purges queues on given RabbitMQ server",
					Action:      purgeCmd,
					Usage: "\n\t" +
						"rabbitr queues purge -s my-server-from-cfg -f \"1==1\"\t# will purge all queues from my-server-from-cfg\n\t" +
						"rabbitr queues purge -s my-server-from-cfg -f 'queue.Name=~\"my-queue\"'\t# will purge queues from my-server-from-cfg which's name matches the 'my-queue' string\n\t" +
						"rabbitr queues purge -s my-server-from-cfg -f 'queue.Messages>200'\t# will purge queues from to my-server-from-cfg with more than 200 messages\n\t" +
						"rabbitr queues purge -s my-server-from-cfg -f 'queue.Consumers==0'\t# will purge queues from to my-server-from-cfg that have 0 consumers\n\t",
				},
			},
		},
	}
}
