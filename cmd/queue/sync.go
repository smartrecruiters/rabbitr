package queue

import (
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func syncCmd(ctx *cli.Context) error {
	executeQueueOperation(ctx, syncQueueFn, nil)
	return nil
}

func syncQueueFn(client *rabbithole.Client, queue *interface{}, w *tabwriter.Writer) {
	q := (*queue).(rabbithole.QueueInfo)
	fmt.Fprintf(w, "Syncing queue: %s/%s\t", q.Vhost, q.Name)
	res, err := client.SyncQueue(q.Vhost, q.Name)
	commons.PrintIfErrorWithMsg(fmt.Sprintf("Error occured when attempting to sync queue %s/%s", q.Vhost, q.Name), err)
	if res != nil {
		fmt.Fprintf(w, "Response code: %d\t", res.StatusCode)
		commons.PrintResponseBodyIfError(res)
	}
}
