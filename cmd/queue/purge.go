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
	commons.Fprintf(w, "Purging queue: %s/%s from %d messages\t", q.Vhost, q.Name, q.Messages)
	res, err := client.PurgeQueue(q.Vhost, q.Name)
	commons.PrintToWriterIfErrorWithMsg(w, fmt.Sprintf("Error occurred when attempting to purge queue %s/%s", q.Vhost, q.Name), err)
	commons.HandleGeneralResponseWithWriter(w, res)
}
