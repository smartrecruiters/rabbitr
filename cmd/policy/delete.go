package policy

import (
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func deleteCmd(ctx *cli.Context) error {
	executePolicyOperation(ctx, deletePolicyFn, nil)
	return nil
}

func deletePolicyFn(client *rabbithole.Client, policy *interface{}, w *tabwriter.Writer) {
	p := (*policy).(rabbithole.Policy)
	commons.Fprintf(w, "Deleting policy: %s/%s \t", p.Vhost, p.Name)
	res, err := client.DeletePolicy(p.Vhost, p.Name)
	commons.PrintToWriterIfErrorWithMsg(w, fmt.Sprintf("Error occurred when attempting to delete queue %s/%s", p.Vhost, p.Name), err)
	commons.HandleGeneralResponseWithWriter(w, res)
}
