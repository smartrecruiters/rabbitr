package connection

import (
	"fmt"
	"text/tabwriter"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/urfave/cli"
)

func listConnectionsCmd(ctx *cli.Context) error {
	executeConnectionOperation(ctx, listConnectionFn, printListConnectionsHeaderFn)
	return nil
}

func printListConnectionsHeaderFn(w *tabwriter.Writer) {
	_, _ = fmt.Fprintln(w, "Connection Name\tClient Provided Name\t")
}

func listConnectionFn(client *rabbithole.Client, connection *interface{}, w *tabwriter.Writer) {
	c := (*connection).(ConnectionInfo)
	fmt.Fprintf(w, "%s/%s \t%s\t", c.Vhost, c.ID, c.Name)
}
