package message

import (
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

// GetCommands returns slice of commands for this command category
func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "message",
			Aliases:     []string{"messages"},
			Hidden:      false,
			Description: "Group of commands related to messages",
			Subcommands: []cli.Command{
				{
					Name: "move",
					Flags: []cli.Flag{
						commons.ServerFlag,
						cli.StringFlag{
							Name:  "src-vhost",
							Usage: "Source vhost",
						},
						cli.StringFlag{
							Name:  "src-queue",
							Usage: "Source queue",
						},
						cli.StringFlag{
							Name:  "dst-vhost",
							Usage: "Optional. Destination vhost, if not provided defaults to vhost",
						},
						cli.StringFlag{
							Name:  "dst-queue",
							Usage: "Destination queue. Omitted if invoked with a duplicate flag",
						},
						cli.BoolFlag{
							Name:  "duplicate",
							Usage: "Optional. Flag indicating that messages from the source queue should be moved to the two newly created separate queues (duplicating them as a result)",
						},
						cli.IntFlag{
							Name:  "prefetch-count",
							Usage: "Optional. Maximum number of unacknowledged messages that may be in flight over a shovel at one time. Defaults to 1000 if not set.",
							Value: 1000,
						},
						cli.StringFlag{
							Name:  "delete-after",
							Usage: "Optional. Determines when (if ever) the underlying shovel should delete itself. Can be one of: [never/queue-length/or a fixed integer number of transferred messages]. When setting it to integer, the underlying shovel will delete itself after transferring that many messages.",
							Value: "10000",
						},
					},
					Description: "Moves messages between queues, it uses shovel under the hood. It can also move messages from a source queue to the two separate queues (duplicating messages as a result).\n\t " +
						"Please take a look at the https://www.rabbitmq.com/shovel-dynamic.html#amqp091-reference for broader description.\n\t " +
						"If shovel created by this command already exist, it will be updated with the consecutive command invocations.",
					Action: moveMessagesCmd,
					Usage: "\n\t" +
						"rabbitr messages move -s my-server-from-cfg -src-vhost test -src-queue my-queue -dst-vhost test2 -dst-queue my-new-dest-queue\t# will move messages from source queue to destination queue on given vhosts\n\t" +
						"rabbitr messages move -s my-server-from-cfg -src-vhost test -src-queue my-queue -dst-queue my-new-dest-queue -delete-after 10000\t# will move 10000 messages from source queue to destination queue\n\t" +
						"rabbitr messages move -s my-server-from-cfg -src-vhost test -src-queue my-queue -dst-queue my-new-dest-queue -delete-after never\t# will create shovel for moving messages from source queue to destination queue, each time message will appear in the source queue it will be moved to the destination one\n\t" +
						"rabbitr messages move -s my-server-from-cfg -src-vhost test -src-queue my-queue -dst-queue my-new-dest-queue -delete-after queue-length\t# will move all messages from source queue to destination queue (avoid on large queues)\n\t" +
						"rabbitr messages move -s my-server-from-cfg -src-vhost test -src-queue my-queue -duplicate \t# will duplicate and move messages from source queue to the two newly created queues\n\t",
				},
			},
		},
	}
}
