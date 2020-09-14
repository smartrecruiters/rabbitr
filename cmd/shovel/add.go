package shovel

import (
	"fmt"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func addShovelCmd(ctx *cli.Context) {
	server := ctx.String("server-name")
	vhost := ctx.String("vhost")
	shovel := ctx.String("shovel")
	client := commons.GetRabbitClient(server)
	definition := rabbithole.ShovelDefinition{
		SourceURI:              ctx.String("src-uri"),
		SourceExchange:         ctx.String("src-exchange"),
		SourceExchangeKey:      ctx.String("src-exchange-key"),
		SourceQueue:            ctx.String("src-queue"),
		DestinationURI:         ctx.String("dst-uri"),
		DestinationExchange:    ctx.String("dst-exchange"),
		DestinationExchangeKey: ctx.String("dst-exchange-key"),
		DestinationQueue:       ctx.String("dst-queue"),
		PrefetchCount:          ctx.Int("prefetch-count"),
		ReconnectDelay:         ctx.Int("reconnect-delay"),
		AddForwardHeaders:      ctx.Bool("add-forward-headers"),
		AckMode:                ctx.String("ack-mode"),
		DeleteAfter:            ctx.String("delete-after"),
	}
	res, err := client.DeclareShovel(vhost, shovel, definition)
	fmt.Printf("%v %s", res, err)
}
