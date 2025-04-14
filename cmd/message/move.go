package message

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/smartrecruiters/rabbitr/cmd/rabbit"

	"github.com/google/uuid"

	"github.com/smartrecruiters/rabbitr/cmd/server"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func moveMessagesCmd(ctx *cli.Context) error {
	s := server.AskForServerSelection(ctx.String(commons.ServerName))
	srcVHost := commons.AskIfValueEmpty(strings.TrimSpace(ctx.String("src-vhost")), "src-vhost")
	srcQueue := commons.AskIfValueEmpty(strings.TrimSpace(ctx.String("src-queue")), "src-queue")
	dstVHost := getDstVHost(ctx.String("dst-vhost"), srcVHost)
	duplicate := ctx.Bool("duplicate")
	client := rabbit.GetRabbitClient(s)

	if duplicate {
		return moveAndDuplicate(client, srcVHost, srcQueue, dstVHost, ctx)
	}
	return moveOnly(client, srcVHost, srcQueue, dstVHost, ctx)
}

func moveOnly(client *rabbithole.Client, srcVHost string, srcQueue string, dstVHost string, ctx *cli.Context) error {
	dstQueue := commons.AskIfValueEmpty(strings.TrimSpace(ctx.String("dst-queue")), "dst-queue")
	name := "Move from " + srcQueue
	declareShovel(client, srcVHost, srcQueue, dstVHost, "", dstQueue, name, ctx.Int("prefetch-count"), getDeleteAfter(ctx))
	return nil
}

func getDeleteAfter(ctx *cli.Context) rabbithole.DeleteAfter {
	deleteAfterStr := ctx.String("delete-after")
	deleteAfterInt, err := strconv.Atoi(deleteAfterStr)
	if err == nil && deleteAfterInt > 0 {
		return rabbithole.DeleteAfter(strconv.Itoa(deleteAfterInt))
	}

	if len(deleteAfterStr) > 0 {
		return rabbithole.DeleteAfter(deleteAfterStr)
	}

	return "queue-length"
}

func getDstVHost(dstVHost, defaultValue string) string {
	dstVHost = strings.TrimSpace(dstVHost)
	if len(dstVHost) <= 0 {
		dstVHost = defaultValue
	}
	return dstVHost
}

func moveAndDuplicate(client *rabbithole.Client, srcVHost, srcQueue, dstVHost string, ctx *cli.Context) error {
	id := uuid.New().String()
	dstQueue1 := srcQueue + ".Mirror1." + id
	dstQueue2 := srcQueue + ".Mirror2." + id
	dstExchange := "RabbitrDuplicatingExchange." + id
	shovelName := "RabbitrDuplicatingShovel." + id

	declareExchange(client, dstVHost, dstExchange)
	declareQueue(client, dstVHost, dstQueue1)
	declareQueue(client, dstVHost, dstQueue2)
	declareBindingForQueue(client, dstVHost, dstExchange, dstQueue1)
	declareBindingForQueue(client, dstVHost, dstExchange, dstQueue2)
	declareShovel(client, srcVHost, srcQueue, dstVHost, dstExchange, "", shovelName, ctx.Int("prefetch-count"), getDeleteAfter(ctx))
	fmt.Printf("Operation completed, please check contents in both %s and %s queues\n", dstQueue1, dstQueue2)
	return nil
}

func declareExchange(client *rabbithole.Client, srcVHost string, dstExchange string) {
	res, err := client.DeclareExchange(srcVHost, dstExchange, rabbithole.ExchangeSettings{
		Type:       "fanout",
		Durable:    false,
		AutoDelete: true,
	})
	commons.HandleGeneralResponse(fmt.Sprintf("Creating %s/%s exchange", srcVHost, dstExchange), res)
	commons.AbortIfErrorWithMsg(fmt.Sprintf("Unable to create %s/%s exchange", srcVHost, dstExchange), err)
}

func declareQueue(client *rabbithole.Client, srcVHost string, queueName string) {
	res, err := client.DeclareQueue(srcVHost, queueName, rabbithole.QueueSettings{
		Durable:    true,
		AutoDelete: false,
		Arguments:  nil,
	})
	commons.HandleGeneralResponse(fmt.Sprintf("Creating %s/%s queue", srcVHost, queueName), res)
	commons.AbortIfErrorWithMsg(fmt.Sprintf("Unable to create queue %s/%s queue", srcVHost, queueName), err)
}

func declareBindingForQueue(client *rabbithole.Client, srcVHost string, dstExchange string, queueName string) {
	res, err := client.DeclareBinding(srcVHost, rabbithole.BindingInfo{
		Source:          dstExchange,
		Vhost:           srcVHost,
		Destination:     queueName,
		DestinationType: "queue",
	})
	commons.HandleGeneralResponse(fmt.Sprintf("Creating binding for %s/%s queue", srcVHost, queueName), res)
	commons.AbortIfErrorWithMsg(fmt.Sprintf("Unable to create binding between exchange %s/%s and queue %s/%s queue", srcVHost, dstExchange, srcVHost, queueName), err)
}

func declareShovel(client *rabbithole.Client, srcVHost, srcQueue, dstVHost, dstExchange, dstQueue, shovelName string, prefetchCount int, deleteAfter rabbithole.DeleteAfter) {
	res, err := client.DeclareShovel(dstVHost, shovelName, rabbithole.ShovelDefinition{
		SourceURI:           rabbithole.URISet{"amqp:///" + srcVHost},
		SourceQueue:         srcQueue,
		DestinationURI:      rabbithole.URISet{"amqp:///" + dstVHost},
		ReconnectDelay:      15,
		DestinationExchange: dstExchange,
		DestinationQueue:    dstQueue,
		PrefetchCount:       prefetchCount,
		AddForwardHeaders:   true,
		AckMode:             "on-confirm",
		DeleteAfter:         deleteAfter,
	})
	commons.HandleGeneralResponse(fmt.Sprintf("Creating temporary shovel %s/%s to move messages", dstVHost, shovelName), res)
	commons.AbortIfError(err)
}
