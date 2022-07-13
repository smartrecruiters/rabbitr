package message

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/smartrecruiters/rabbitr/cmd/rabbit"
	"github.com/smartrecruiters/rabbitr/cmd/server"
	"github.com/urfave/cli"
)

func downloadMessagesCmd(ctx *cli.Context) error {
	s := server.AskForServerSelection(ctx.String(commons.ServerName))
	vHost := commons.AskIfValueEmpty(strings.TrimSpace(ctx.String("vhost")), "vhost")
	queue := commons.AskIfValueEmpty(strings.TrimSpace(ctx.String("queue")), "queue")
	maxMessages := ctx.Int("max-messages")
	prefetchCount := ctx.Int("prefetch-count")
	outDir := ctx.String("output-dir")
	verbose := ctx.Bool("verbose")

	consumer := rabbit.InitConsumer()
	consumer.Connection = rabbit.GetAmqpRabbitConnection(s, vHost)
	subscribeConsumer(consumer, queue, prefetchCount, outDir, maxMessages, verbose)

	commons.PrintIfTrue(verbose, fmt.Sprintf("Running until all %d message(s) are consumed (to abort early press Ctrl+C)", maxMessages))
	<-consumer.Done
	commons.PrintIfTrue(verbose, "Shutting down...")
	err := consumer.Shutdown(verbose)
	commons.DebugIfError(err)
	return nil
}

func subscribeConsumer(consumer *rabbit.Consumer, srcQueue string, prefetchCount int, outDir string, maxMessages int, verbose bool) {
	var err error
	commons.Debug("Got connection, creating channel")
	consumer.Channel, err = consumer.Connection.Channel()
	commons.AbortIfError(err)
	err = consumer.Channel.Qos(prefetchCount, 0, false)
	commons.AbortIfError(err)

	commons.Debug("Starting consumption")
	deliveries, err := consumer.Channel.Consume(
		srcQueue,
		consumer.Tag,
		false,
		false,
		false,
		false,
		nil,
	)
	commons.AbortIfError(err)
	rabbit.SetupCloseHandler(consumer, verbose)
	err = commons.MakeDir(outDir)
	commons.AbortIfError(err)

	go handleDelivery(deliveries, consumer.Done, outDir, maxMessages, verbose)
}

func handleDelivery(deliveries <-chan amqp.Delivery, done chan error, outputDir string, maxMessagesToConsume int, verbose bool) {
	cleanup := func() {
		commons.Debug("HandleDelivery: deliveries channel closed")
		done <- nil
	}

	defer cleanup()

	messagesReceived := 0
	for msg := range deliveries {
		messagesReceived++

		err := saveMessageToFile(msg.Body, outputDir, messagesReceived)
		commons.AbortIfError(err)

		err = savePropsAndHeadersToFile(msg, outputDir, messagesReceived)
		commons.AbortIfError(err)

		err = msg.Ack(false)
		commons.AbortIfError(err)
		if messagesReceived == maxMessagesToConsume {
			commons.PrintIfTrue(verbose, fmt.Sprintf("Received %d message(s), finishing up", messagesReceived))
			done <- nil
			break
		}
	}
}

func saveMessageToFile(body []byte, outputDir string, counter int) error {
	filePath := generateFilePath(outputDir, counter, "")
	err := ioutil.WriteFile(filePath, body, 0644)
	if err != nil {
		return err
	}

	fmt.Println(filePath)

	return nil
}

func savePropsAndHeadersToFile(msg amqp.Delivery, outputDir string, counter int) error {
	extras := rabbit.InitMessageExtras(msg)

	data, err := json.MarshalIndent(extras, "", "  ")
	if err != nil {
		return err
	}

	filePath := generateFilePath(outputDir, counter, "-headers+properties.json")
	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	fmt.Println(filePath)

	return nil
}

func generateFilePath(outputDir string, counter int, optionalSuffix string) string {
	return path.Join(outputDir, fmt.Sprintf("msg-%04d%s", counter, optionalSuffix))
}
