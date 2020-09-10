package connection

import (
	"fmt"
	"text/tabwriter"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/urfave/cli"
)

func closeConnectionsCmd(ctx *cli.Context) error {
	executeConnectionOperation(ctx, closeConnectionFn, nil)
	return nil
}

func closeConnectionFn(client *rabbithole.Client, connection *interface{}, w *tabwriter.Writer) {
	c := (*connection).(ConnectionInfo)
	fmt.Fprintf(w, "Closing connection %s/%s with name: %s\t", c.Vhost, c.ID, c.Name)
	res, err := client.CloseConnection(c.ID)
	commons.PrintIfErrorWithMsg(fmt.Sprintf("Error occured when attempting to close connection %s/%s", c.Vhost, c.ID), err)
	if res != nil {
		fmt.Fprintf(w, "Response code: %d\t", res.StatusCode)
		commons.PrintResponseBodyIfError(res)
	}
}
