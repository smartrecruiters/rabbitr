package queue

import (
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func purgeCmd(ctx *cli.Context) error {
	executeQueueOperation(ctx, purgeQueueFn, nil)
	return nil
}

func purgeQueueFn(client *rabbithole.Client, queue *interface{}, w *tabwriter.Writer) {
	q := (*queue).(rabbithole.QueueInfo)
	fmt.Fprintf(w, "Purging queue: %s/%s from %d messages\t", q.Vhost, q.Name, q.Messages)
	res, err := client.PurgeQueue(q.Vhost, q.Name)
	commons.PrintIfErrorWithMsg(fmt.Sprintf("Error occured when attempting to purge queue %s/%s", q.Vhost, q.Name), err)
	if res != nil {
		fmt.Fprintf(w, "Response code: %d\t", res.StatusCode)
		commons.PrintResponseBodyIfError(res)
	}
}
