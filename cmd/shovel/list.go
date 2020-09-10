package shovel

import (
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole"
)

func listShovelsCmd(ctx *cli.Context) error {
	executeShovelOperation(ctx, listShovelFn, printListShovelsHeaderFn)
	return nil
}

func printListShovelsHeaderFn(w *tabwriter.Writer) {
	_, _ = fmt.Fprintln(w, "Shovel\tDefinition\t")
}

func listShovelFn(client *rabbithole.Client, shovel *interface{}, w *tabwriter.Writer) {
	s := (*shovel).(rabbithole.ShovelInfo)
	fmt.Fprintf(w, "%s/%s \t%v\t", s.Vhost, s.Name, s.Definition)
}
