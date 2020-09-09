package queue

import (
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole"
)

func listQueuesCmd(ctx *cli.Context) error {
	executeQueueOperation(ctx, listQueueFn, printListQueuesHeaderFn)
	return nil
}

func printListQueuesHeaderFn(w *tabwriter.Writer) {
	_, _ = fmt.Fprintln(w, "Queue\tConsumers\tMessages\tOwner\t")
}

func listQueueFn(client *rabbithole.Client, queue *interface{}, w *tabwriter.Writer) {
	q := (*queue).(rabbithole.QueueInfo)
	owner := getQueueOwner(q)
	fmt.Fprintf(w, "%s/%s \t%d\t%d\t%s\t", q.Vhost, q.Name, q.Consumers, q.Messages, owner)
}

func getQueueOwner(q rabbithole.QueueInfo) string {
	owner := ""
	ownerArg := q.Arguments["x-queue-owner"]
	if ownerArg != nil {
		owner = ownerArg.(string)
	}
	return owner
}
