package exchange

import (
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func deleteCmd(ctx *cli.Context) error {
	executeExchangeOperation(ctx, deleteExchangeFn, nil)
	return nil
}

func deleteExchangeFn(client *rabbithole.Client, exchange *interface{}, w *tabwriter.Writer) {
	e := (*exchange).(rabbithole.ExchangeInfo)
	fmt.Fprintf(w, "Deleting exchange: %s/%s \t", e.Vhost, e.Name)
	res, err := client.DeleteExchange(e.Vhost, e.Name)
	commons.PrintIfErrorWithMsg(fmt.Sprintf("Error occured when attempting to delete an exchange %s/%s", e.Vhost, e.Name), err)
	if res != nil {
		fmt.Fprintf(w, "Response code: %d\t", res.StatusCode)
		commons.PrintResponseBodyIfError(res)
	}
}
