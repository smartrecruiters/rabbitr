package policy

import (
	"text/tabwriter"

	"github.com/smartrecruiters/rabbitr/cmd/commons"

	"github.com/urfave/cli"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
)

func printListPoliciesHeaderFn(w *tabwriter.Writer) {
	commons.Fprintln(w, "Policy\tPattern\tPriority\tApplyTo\tDefinition\t")
}

func showPolicyFn(client *rabbithole.Client, policy *interface{}, w *tabwriter.Writer) {
	p := (*policy).(rabbithole.Policy)
	commons.Fprintf(w, "%s/%s\t%s\t%d\t%s\t%v\t", p.Vhost, p.Name, p.Pattern, p.Priority, p.ApplyTo, p.Definition)
}

func listPoliciesCmd(ctx *cli.Context) error {
	executePolicyOperation(ctx, showPolicyFn, printListPoliciesHeaderFn)
	return nil
}
