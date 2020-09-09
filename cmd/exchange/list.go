package exchange

import (
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole"
)

func listExchangesCmd(ctx *cli.Context) error {
	executeExchangeOperation(ctx, listExchangeFn, printListExchangesHeaderFn)
	return nil
}

func printListExchangesHeaderFn(w *tabwriter.Writer) {
	_, _ = fmt.Fprintln(w, "Exchange\tDurable\tAutoDelete\tInternal\t")
}

func listExchangeFn(client *rabbithole.Client, exchange *interface{}, w *tabwriter.Writer) {
	e := (*exchange).(rabbithole.ExchangeInfo)
	fmt.Fprintf(w, "%s/%s \t%v\t%v\t%v\t", e.Vhost, e.Name, e.Durable, e.AutoDelete, e.Internal)
}
