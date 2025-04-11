package shovel

import (
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func deleteCmd(ctx *cli.Context) error {
	executeShovelOperation(ctx, deleteShovelFn, nil)
	return nil
}

func deleteShovelFn(client *rabbithole.Client, shovel *interface{}, w *tabwriter.Writer) {
	s := (*shovel).(rabbithole.ShovelInfo)
	commons.Fprintf(w, "Deleting shovel: %s/%s \t", s.Vhost, s.Name)
	res, err := client.DeleteShovel(s.Vhost, s.Name)
	commons.PrintToWriterIfErrorWithMsg(w, fmt.Sprintf("Error occurred when attempting to delete shovel %s/%s", s.Vhost, s.Name), err)
	commons.HandleGeneralResponseWithWriter(w, res)
}
