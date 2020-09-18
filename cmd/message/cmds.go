package message

import (
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

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
					},
					Description: "Moves messages between queues, it uses shovel under the hood. It can also move messages from a source queue to the two separate queues (duplicating them as a result)",
					Action:      moveMessagesCmd,
					Usage: "\n\t" +
						"rabbitr messages move -s my-server-from-cfg -src-vhost test -src-queue my-queue -dst-vhost test2 -dst-queue my-new-dest-queue\t# will move messages from source queue to destination queue on given vhosts\n\t" +
						"rabbitr messages move -s my-server-from-cfg -src-vhost test -src-queue my-queue -duplicate \t# will duplicate and move messages from source queue to the two newly created queues\n\t",
				},
			},
		},
	}
}
