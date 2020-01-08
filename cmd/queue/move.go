package queue

import (
	"fmt"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func moveMessagesCmd(ctx *cli.Context) error {
	server := ctx.String("server-name")
	srcVHost := ctx.String("src-vhost")
	dstVHost := ctx.String("dst-vhost")
	srcQueue := ctx.String("src-queue")
	dstQueue := ctx.String("dst-queue")
	name := "Move from " + srcQueue
	if len(dstVHost) <= 0 {
		dstVHost = srcVHost
	}

	definition := rabbithole.ShovelDefinition{
		SourceURI:         "amqp:///" + srcVHost,
		SourceQueue:       srcQueue,
		DestinationURI:    "amqp:///" + dstVHost,
		DestinationQueue:  dstQueue,
		PrefetchCount:     1000,
		AddForwardHeaders: false,
		AckMode:           "on-confirm",
		DeleteAfter:       "queue-length",
	}

	client := commons.GetRabbitClient(server)
	res, err := client.DeclareShovel(srcVHost, name, definition)
	commons.AbortIfError(err)
	fmt.Printf("Created temporary shovel to move messages, Response code: %d\t\n", res.StatusCode)
	return nil
}
