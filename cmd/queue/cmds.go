package queue

import (
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "queue",
			Aliases: []string{"queues"},
			Hidden:  false,
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
					Name: "duplicate",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.VHostFlag,
						cli.StringFlag{
							Name:  "queue, q",
							Usage: "Source queue",
						},
					},
					Description: "Moves messages from a source queue to the two separate queues. Messages will be duplicated in newly created queues and removed from the source queue.",
					Action:      duplicateCmd,
					Usage: "\n\t" +
						"rabbitr queues duplicate -s my-server-from-cfg -q my-queue \t# will duplicate and move messages from source queue to the two newly created queues\n\t",
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
				{
					Name: "sync",
					Flags: []cli.Flag{
						commons.ServerFlag,
						commons.DryRunFlag,
						commons.VHostFlag,
						commons.GetFilterFlag(commons.AllSubjects, commons.QueueFilterFields),
					},
					Description: "Sync queues on given RabbitMQ server",
					Action:      syncCmd,
					Usage: "\n\t" +
						"rabbitr queues sync -s my-server-from-cfg \t# will sync all queues from my-server-from-cfg\n\t" +
						"rabbitr queues sync -s my-server-from-cfg -f 'queue.Name=~\"my-queue\"'\t# will sync queues from my-server-from-cfg which's name matches the 'my-queue' string\n\t" +
						"rabbitr queues sync -s my-server-from-cfg -f 'queue.Messages>200'\t# will sync queues from to my-server-from-cfg with more than 200 messages\n\t" +
						"rabbitr queues sync -s my-server-from-cfg -f 'queue.Consumers==0'\t# will sync queues from to my-server-from-cfg that have 0 consumers\n\t",
				},
				{
					Name:    "move-messages",
					Aliases: []string{"move"},
					Flags: []cli.Flag{
						commons.ServerFlag,
						cli.StringFlag{
							Name:  "src-vhost",
							Usage: "Source vhost",
						},
						cli.StringFlag{
							Name:  "dst-vhost",
							Usage: "Optional. Destination vhost, if not provided defaults to vhost",
						},
						cli.StringFlag{
							Name:  "src-queue",
							Usage: "Source queue",
						},
						cli.StringFlag{
							Name:  "dst-queue",
							Usage: "Destination queue",
						},
					},
					Description: "Moves messages between queues, it uses shovel under the hood",
					Action:      moveMessagesCmd,
					Usage: "\n\t" +
						"rabbitr queue move-messages -s my-server-from-cfg -src-vhost test -src-queue my-queue -dst-vhost test2 -dst-queue my-new-dest-queue\t# will move messages from source queue to destination queue on given vhosts\n\t",
				},
			},
		},
	}
}
