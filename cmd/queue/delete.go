package queue

import (
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func deleteCmd(ctx *cli.Context) error {
	executeQueueOperation(ctx, deleteQueueFn, nil)
	return nil
}

func deleteQueueFn(client *rabbithole.Client, queue *interface{}, w *tabwriter.Writer) {
	q := (*queue).(rabbithole.QueueInfo)
	fmt.Fprintf(w, "Deleting queue: %s/%s \t", q.Vhost, q.Name)
	res, err := client.DeleteQueue(q.Vhost, q.Name)
	commons.PrintIfErrorWithMsg(fmt.Sprintf("Error occured when attempting to delete a queue %s/%s", q.Vhost, q.Name), err)
	if res != nil {
		fmt.Fprintf(w, "Response code: %d\t", res.StatusCode)
		commons.PrintResponseBodyIfError(res)
	}
}
