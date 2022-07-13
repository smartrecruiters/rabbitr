package message

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/smartrecruiters/rabbitr/cmd/rabbit"
	"github.com/smartrecruiters/rabbitr/cmd/server"
	"github.com/urfave/cli"
)

var (
	publishSequenceToFileNameMap map[uint64]string
)

func uploadMessagesCmd(ctx *cli.Context) error {
	var err error
	s := server.AskForServerSelection(ctx.String(commons.ServerName))
	vHost := commons.AskIfValueEmpty(strings.TrimSpace(ctx.String("vhost")), "vhost")
	queue := commons.AskIfValueEmpty(strings.TrimSpace(ctx.String("queue")), "queue")
	maxMessages := commons.AskForIntIf(commons.LessOrEqualThanZeroValidator, ctx.Int("max-messages"), "max-messages")
	inputDir := ctx.String("input-dir")
	verbose := ctx.Bool("verbose")
	confirm := ctx.Bool("confirm")
	predeclared := ctx.Bool("predeclared")
	removeOnConfirm := ctx.Bool("remove-on-confirm")

	publisher := rabbit.InitPublisher()
	publisher.Connection = rabbit.GetAmqpRabbitConnection(s, vHost)
	publisher.Channel, err = publisher.Connection.Channel()
	publisher.MessagesToPublish = maxMessages
	publishSequenceToFileNameMap = make(map[uint64]string, maxMessages)

	if confirm {
		setupPublisherConfirms(publisher, inputDir, removeOnConfirm, verbose)
	}

	ensureQueueExists(publisher, queue, predeclared)
	files, err := ioutil.ReadDir(inputDir)
	commons.AbortIfError(err)

	for _, file := range files {
		fileName := file.Name()
		// skip directories, files that do not start with 'msg-' and files 'msg-*.json'
		if file.IsDir() || !strings.HasPrefix(fileName, "msg-") || strings.Contains(fileName, ".") {
			continue
		}

		seqNo := publisher.Channel.GetNextPublishSeqNo()
		registerPublishedMessage(seqNo, fileName)
		msg, extras := readMessageAndHeaders(inputDir, fileName)
		exchange := commons.GetStringOrDefaultIfValue(queue, "", extras.Properties.Exchange)
		routingKey := commons.GetStringOrDefault(queue, extras.Properties.RoutingKey)
		err = publisher.Channel.Publish(exchange, routingKey, false, false, msg)
		commons.AbortIfErrorWithMsg("Unable to publish message due to: %v", err)
		if confirm {
			publisher.Publishes <- seqNo
		}

		fmt.Println(getFilePath(inputDir, fileName))
		fmt.Println(getHeadersFilePath(inputDir, fileName))
		publisher.NumberOfPublished++

		if publisher.NumberOfPublished >= maxMessages {
			commons.PrintIfTrue(verbose, "All messages published.")
			publisher.AllPublished = true
			break
		}
	}

	if confirm {
		commons.PrintIfTrue(verbose, "Awaiting for all confirmations...")
		<-publisher.AllPublishedAndConfirmed
	}

	commons.PrintIfTrue(verbose, fmt.Sprintf("Published %d messages, shutting down...", publisher.NumberOfPublished))
	err = publisher.Shutdown(verbose)
	commons.DebugIfError(err)
	return nil
}

func ensureQueueExists(publisher *rabbit.Publisher, queue string, predeclared bool) {
	channel, err := publisher.Connection.Channel()
	commons.AbortIfError(err)
	_, err = channel.QueueDeclare(queue, true, false, false, false, amqp.Table{})
	commons.AbortIfTrue(err != nil && !predeclared, fmt.Sprintf("Unable to create queue, use --predeclared if a queue with diffrent configuration already exists, err: %s", err))
}

func registerPublishedMessage(seqNo uint64, fileName string) {
	publishSequenceToFileNameMap[seqNo] = fileName
}

func setupPublisherConfirms(publisher *rabbit.Publisher, inputDir string, removeOnConfirm bool, verbose bool) {
	commons.Debug("Enabling publisher confirms.")
	err := publisher.Channel.Confirm(false)
	commons.AbortIfErrorWithMsg("Channel could not be put into confirm mode: %s", err)

	publisher.Confirms = publisher.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	go confirmHandler(publisher, inputDir, removeOnConfirm, verbose)
}

func readMessageAndHeaders(inputDir, fileName string) (amqp.Publishing, rabbit.MessageExtras) {
	msgBody, msgBodyReadErr := ioutil.ReadFile(getFilePath(inputDir, fileName))
	commons.AbortIfError(msgBodyReadErr)

	fileWithHeaders := getHeadersFilePath(inputDir, fileName)
	if !commons.FileExists(fileWithHeaders) {
		commons.AbortIfError(fmt.Errorf("file %s is missing, aborting publishing", fileWithHeaders))
	}
	msgExtras, msgExtrasReadErr := ioutil.ReadFile(fileWithHeaders)
	commons.AbortIfError(msgExtrasReadErr)
	extras := rabbit.MessageExtras{}
	err := json.Unmarshal(msgExtras, &extras)
	commons.AbortIfError(err)

	msg := rabbit.CreateMessageToPublish(msgBody, extras)
	return msg, extras
}

func getHeadersFilePath(inputDir, fileName string) string {
	return getFilePath(inputDir, fileName+"-headers+properties.json")
}

func getFilePath(dir, fileName string) string {
	return path.Join(dir, fileName)
}

func confirmHandler(publisher *rabbit.Publisher, inputDir string, removeOnConfirm, verbose bool) {
	confirmations := make(map[uint64]bool)
	for {
		select {
		case publishSeqNo := <-publisher.Publishes:
			commons.PrintIfTrue(verbose, fmt.Sprintf("Waiting for confirmation of %d", publishSeqNo))
			confirmations[publishSeqNo] = false
		case confirmed := <-publisher.Confirms:
			if confirmed.DeliveryTag > 0 {
				if confirmed.Ack {
					commons.PrintIfTrue(verbose, fmt.Sprintf("Confirmed delivery with delivery tag: %d", confirmed.DeliveryTag))
					if removeOnConfirm {
						fileWithMessage := getFilePath(inputDir, publishSequenceToFileNameMap[confirmed.DeliveryTag])
						fileWithHeaders := getHeadersFilePath(inputDir, publishSequenceToFileNameMap[confirmed.DeliveryTag])
						commons.Debugf("Removing file: %s", fileWithMessage)
						err := os.Remove(fileWithMessage)
						commons.DebugIfError(err)
						commons.Debugf("Removing file: %s", fileWithHeaders)
						err = os.Remove(fileWithHeaders)
						commons.DebugIfError(err)
					}
				} else {
					commons.PrintIfTrue(verbose, fmt.Sprintf("Failed delivery of delivery tag: %d", confirmed.DeliveryTag))
				}
				delete(confirmations, confirmed.DeliveryTag)
				publisher.NumberOfConfirmed++
				commons.Debugf("allPublished: %v, NumberOfPublished: %d, NumberOfConfirmed:%d, MessagesToPublish: %d", publisher.AllPublished, publisher.NumberOfPublished, publisher.NumberOfConfirmed, publisher.MessagesToPublish)
				if publisher.AllPublished && publisher.NumberOfPublished == publisher.NumberOfConfirmed && publisher.NumberOfConfirmed == publisher.MessagesToPublish {
					commons.Debug("All published and confirmed, confirmHandler is stopping")
					publisher.AllPublishedAndConfirmed <- true
					return
				}
			}
		}
		if len(confirmations) > 1 {
			commons.PrintIfTrue(verbose, fmt.Sprintf("Outstanding confirmations: %d", len(confirmations)))
		}
	}
}
