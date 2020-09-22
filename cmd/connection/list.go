package connection

import (
	"text/tabwriter"

	"github.com/smartrecruiters/rabbitr/cmd/commons"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/urfave/cli"
)

func listConnectionsCmd(ctx *cli.Context) error {
	executeConnectionOperation(ctx, listConnectionFn, printListConnectionsHeaderFn)
	return nil
}

func printListConnectionsHeaderFn(w *tabwriter.Writer) {
	commons.Fprintln(w, "Connection Name\tClient Provided Name\t")
}

func listConnectionFn(client *rabbithole.Client, connection *interface{}, w *tabwriter.Writer) {
	c := (*connection).(ConnInfo)
	commons.Fprintf(w, "%s/%s \t%s\t", c.Vhost, c.ID, c.Name)
}
