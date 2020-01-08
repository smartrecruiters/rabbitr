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
	_, _ = fmt.Fprintln(w, "Queue\tConsumers\tMessages\t")
}

func listQueueFn(client *rabbithole.Client, queue *interface{}, w *tabwriter.Writer) {
	q := (*queue).(rabbithole.QueueInfo)
	fmt.Fprintf(w, "%s/%s \t%d\t%d\t", q.Vhost, q.Name, q.Consumers, q.Messages)
}
