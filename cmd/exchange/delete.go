package exchange

import (
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func deleteCmd(ctx *cli.Context) error {
	executeExchangeOperation(ctx, deleteExchangeFn, nil)
	return nil
}

func deleteExchangeFn(client *rabbithole.Client, exchange *interface{}, w *tabwriter.Writer) {
	e := (*exchange).(rabbithole.ExchangeInfo)
	commons.Fprintf(w, "Deleting exchange: %s/%s \t", e.Vhost, e.Name)
	res, err := client.DeleteExchange(e.Vhost, e.Name)
	commons.PrintToWriterIfErrorWithMsg(w, fmt.Sprintf("Error occured when attempting to delete an exchange %s/%s", e.Vhost, e.Name), err)
	commons.HandleGeneralResponseWithWriter(w, res)
}
