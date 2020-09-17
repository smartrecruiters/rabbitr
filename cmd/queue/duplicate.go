package queue

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/smartrecruiters/rabbitr/cmd/server"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func duplicateCmd(ctx *cli.Context) error {
	s := server.AskForServerSelection(ctx.String(commons.ServerName))
	srcVHost := strings.TrimSpace(ctx.String(commons.VHost))
	srcQueue := commons.AskIfValueEmpty(strings.TrimSpace(ctx.String("queue")), "queue")

	id := uuid.New().String()
	dstQueue1 := srcQueue + ".Mirror1." + id
	dstQueue2 := srcQueue + ".Mirror2." + id
	dstExchange := "RabbitrDuplicatingExchange." + id
	shovelName := "RabbitrDuplicatingShovel." + id

	client := commons.GetRabbitClient(s)
	declareExchange(client, srcVHost, dstExchange)
	declareQueue(client, srcVHost, dstQueue1)
	declareQueue(client, srcVHost, dstQueue2)
	declareBindingForQueue(client, srcVHost, dstExchange, dstQueue1)
	declareBindingForQueue(client, srcVHost, dstExchange, dstQueue2)
	declareShovel(client, srcVHost, srcQueue, dstExchange, shovelName)
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

func declareShovel(client *rabbithole.Client, srcVHost string, srcQueue string, dstExchange string, shovelName string) {
	res, err := client.DeclareShovel(srcVHost, shovelName, rabbithole.ShovelDefinition{
		SourceURI:           "amqp:///" + srcVHost,
		SourceQueue:         srcQueue,
		DestinationURI:      "amqp:///" + srcVHost,
		ReconnectDelay:      15,
		DestinationExchange: dstExchange,
		PrefetchCount:       1000,
		AddForwardHeaders:   true,
		AckMode:             "on-confirm",
		DeleteAfter:         "queue-length",
	})
	commons.HandleGeneralResponse(fmt.Sprintf("Creating temporary shovel %s/%s to move messages", srcVHost, shovelName), res)
	commons.AbortIfError(err)
}
