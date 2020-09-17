package queue

import (
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func deleteCmd(ctx *cli.Context) error {
	executeQueueOperation(ctx, deleteQueueFn, nil)
	return nil
}

func deleteQueueFn(client *rabbithole.Client, queue *interface{}, w *tabwriter.Writer) {
	q := (*queue).(rabbithole.QueueInfo)
	commons.Fprintf(w, "Deleting queue: %s/%s \t", q.Vhost, q.Name)
	res, err := client.DeleteQueue(q.Vhost, q.Name)
	commons.PrintToWriterIfErrorWithMsg(w, fmt.Sprintf("Error occured when attempting to delete queue %s/%s", q.Vhost, q.Name), err)
	commons.HandleGeneralResponseWithWriter(w, res)
}
