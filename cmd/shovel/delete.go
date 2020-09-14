package shovel

import (
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func deleteCmd(ctx *cli.Context) error {
	executeShovelOperation(ctx, deleteShovelFn, nil)
	return nil
}

func deleteShovelFn(client *rabbithole.Client, shovel *interface{}, w *tabwriter.Writer) {
	s := (*shovel).(rabbithole.ShovelInfo)
	fmt.Fprintf(w, "Deleting shovel: %s/%s \t", s.Vhost, s.Name)
	res, err := client.DeleteShovel(s.Vhost, s.Name)
	commons.PrintIfErrorWithMsg(fmt.Sprintf("Error occured when attempting to delete a shovel %s/%s", s.Vhost, s.Name), err)
	if res != nil {
		fmt.Fprintf(w, "Response code: %d\t", res.StatusCode)
		commons.PrintResponseBodyIfError(res)
	}
}
