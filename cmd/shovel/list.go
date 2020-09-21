package shovel

import (
	"text/tabwriter"

	"github.com/smartrecruiters/rabbitr/cmd/commons"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func listShovelsCmd(ctx *cli.Context) error {
	executeShovelOperation(ctx, listShovelFn, printListShovelsHeaderFn)
	return nil
}

func printListShovelsHeaderFn(w *tabwriter.Writer) {
	commons.Fprintln(w, "Shovel\tDefinition\t")
}

func listShovelFn(client *rabbithole.Client, shovel *interface{}, w *tabwriter.Writer) {
	s := (*shovel).(rabbithole.ShovelInfo)
	commons.Fprintf(w, "%s/%s \t%+v\t", s.Vhost, s.Name, s.Definition)
}
