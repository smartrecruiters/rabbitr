package queue

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/smartrecruiters/rabbitr/cmd/server"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func duplicateCmd(ctx *cli.Context) error {
	s := server.AskForServerSelection(ctx.String("server-name"))
	srcVHost := strings.TrimSpace(ctx.String("vhost"))
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
	exchangeCreationResp, err := client.DeclareExchange(srcVHost, dstExchange, rabbithole.ExchangeSettings{
		Type:       "fanout",
		Durable:    false,
		AutoDelete: true,
		Arguments:  nil,
	})
	commons.AbortIfErrorWithMsg(fmt.Sprintf("Unable to create %s/%s exchange", srcVHost, dstExchange), err)
	if exchangeCreationResp != nil {
		fmt.Printf("Creating %s/%s exchange, Response code: %d\t\n", srcVHost, dstExchange, exchangeCreationResp.StatusCode)
		commons.PrintResponseBodyIfError(exchangeCreationResp)
	}
}

func declareQueue(client *rabbithole.Client, srcVHost string, queueName string) {
	queueCreationResp, err := client.DeclareQueue(srcVHost, queueName, rabbithole.QueueSettings{
		Type:       "",
		Durable:    true,
		AutoDelete: false,
		Arguments:  nil,
	})
	commons.AbortIfErrorWithMsg(fmt.Sprintf("Unable to create queue %s/%s queue", srcVHost, queueName), err)
	if queueCreationResp != nil {
		fmt.Printf("Creating %s/%s queue, Response code: %d\t\n", srcVHost, queueName, queueCreationResp.StatusCode)
		commons.PrintResponseBodyIfError(queueCreationResp)
	}
}

func declareBindingForQueue(client *rabbithole.Client, srcVHost string, dstExchange string, queueName string) {
	bindingCreationResp, err := client.DeclareBinding(srcVHost, rabbithole.BindingInfo{
		Source:          dstExchange,
		Vhost:           srcVHost,
		Destination:     queueName,
		DestinationType: "queue",
		RoutingKey:      "",
		Arguments:       nil,
		PropertiesKey:   "",
	})
	commons.AbortIfErrorWithMsg(fmt.Sprintf("Unable to create binding between exchange %s/%s and queue %s/%s queue", srcVHost, dstExchange, srcVHost, queueName), err)
	if bindingCreationResp != nil {
		fmt.Printf("Creating binding for %s/%s queue, Response code: %d\t\n", srcVHost, queueName, bindingCreationResp.StatusCode)
		commons.PrintResponseBodyIfError(bindingCreationResp)
	}
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
	if res != nil {
		fmt.Printf("Creating temporary shovel %s/%s to move messages, Response code: %d\t\n", srcVHost, shovelName, res.StatusCode)
		commons.PrintResponseBodyIfError(res)
	}
	commons.AbortIfError(err)
}
