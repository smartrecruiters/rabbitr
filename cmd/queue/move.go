package queue

import (
	"fmt"
	"strings"

	"github.com/smartrecruiters/rabbitr/cmd/server"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func moveMessagesCmd(ctx *cli.Context) error {
	s := server.AskForServerSelection(ctx.String("server-name"))
	srcVHost := commons.AskIfValueEmpty(strings.TrimSpace(ctx.String("src-vhost")), "src-vhost")
	srcQueue := commons.AskIfValueEmpty(strings.TrimSpace(ctx.String("src-queue")), "src-queue")
	dstQueue := commons.AskIfValueEmpty(strings.TrimSpace(ctx.String("dst-queue")), "dst-queue")
	dstVHost := strings.TrimSpace(ctx.String("dst-vhost"))

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

	client := commons.GetRabbitClient(s)
	res, err := client.DeclareShovel(srcVHost, name, definition)
	if res != nil {
		fmt.Printf("Creating temporary shovel to move messages, Response code: %d\t\n", res.StatusCode)
		commons.PrintResponseBodyIfError(res)
	}
	commons.AbortIfError(err)
	return nil
}
