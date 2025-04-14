package queue

import (
	"text/tabwriter"

	"github.com/smartrecruiters/rabbitr/cmd/commons"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
)

func listQueuesCmd(ctx *cli.Context) error {
	executeQueueOperation(ctx, listQueueFn, printListQueuesHeaderFn)
	return nil
}

func printListQueuesHeaderFn(w *tabwriter.Writer) {
	commons.Fprintln(w, "Queue\tConsumers\tMessages\tOwner\t")
}

func listQueueFn(client *rabbithole.Client, queue *interface{}, w *tabwriter.Writer) {
	q := (*queue).(rabbithole.QueueInfo)
	owner := getQueueOwner(q)
	commons.Fprintf(w, "%s/%s\t%d\t%d\t%s\t", q.Vhost, q.Name, q.Consumers, q.Messages, owner)
}

func getQueueOwner(q rabbithole.QueueInfo) string {
	owner := ""
	ownerArg := q.Arguments["x-queue-owner"]
	if ownerArg != nil {
		owner = ownerArg.(string)
	}
	return owner
}
