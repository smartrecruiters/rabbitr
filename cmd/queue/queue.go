package queue

import (
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func getQueues(client *rabbithole.Client, vhost string) (*[]rabbithole.QueueInfo, error) {
	var queues []rabbithole.QueueInfo
	var err error
	if len(vhost) > 0 {
		queues, err = client.ListQueuesIn(vhost)
	} else {
		queues, err = client.ListQueues()
	}
	commons.AbortIfError(err)
	return &queues, err
}

func getQueueName(subject *interface{}) string {
	q := (*subject).(rabbithole.QueueInfo)
	return q.Name
}

func executeQueueOperation(ctx *cli.Context, queueActionFn commons.SubjectActionFn, printHeaderFn commons.HeaderPrinterFn) {
	server := ctx.String("server-name")
	vhost := ctx.String("vhost")

	client := commons.GetRabbitClient(server)
	queues, err := getQueues(client, vhost)
	commons.AbortIfError(err)

	subjects := commons.ConvertToSliceOfInterfaces(*queues)
	subjectOperator := commons.SubjectOperator{
		ExecuteAction: queueActionFn,
		GetName:       getQueueName,
		Type:          "queue",
		PrintHeader:   printHeaderFn,
	}
	commons.ExecuteOperation(ctx, client, &subjects, subjectOperator)
}
