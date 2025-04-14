package queue

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/smartrecruiters/rabbitr/cmd/rabbit"
	"github.com/smartrecruiters/rabbitr/cmd/server"
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
	return &queues, err
}

func getQueueName(subject *interface{}) string {
	q := (*subject).(rabbithole.QueueInfo)
	return q.Name
}

func executeQueueOperation(ctx *cli.Context, queueActionFn commons.SubjectActionFn, printHeaderFn commons.HeaderPrinterFn) {
	s := server.AskForServerSelection(ctx.String(commons.ServerName))
	vhost := ctx.String(commons.VHost)

	client := rabbit.GetRabbitClient(s)
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
