package shovel

import (
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/smartrecruiters/rabbitr/cmd/server"
	"github.com/urfave/cli"
)

func getPolicies(client *rabbithole.Client, vhost string) (*[]rabbithole.Policy, error) {
	var policies []rabbithole.Policy
	var err error
	if len(vhost) > 0 {
		policies, err = client.ListPoliciesIn(vhost)
	} else {
		policies, err = client.ListPolicies()
	}
	return &policies, err
}

func getPolicyName(subject *interface{}) string {
	p := (*subject).(rabbithole.Policy)
	return p.Name
}

func executePolicyOperation(ctx *cli.Context, policyActionFn commons.SubjectActionFn, printHeaderFn commons.HeaderPrinterFn) {
	s := server.AskForServerSelection(ctx.String("server-name"))
	vhost := ctx.String("vhost")

	client := commons.GetRabbitClient(s)
	policies, err := getPolicies(client, vhost)
	commons.AbortIfError(err)

	subjects := commons.ConvertToSliceOfInterfaces(*policies)
	subjectOperator := commons.SubjectOperator{
		ExecuteAction: policyActionFn,
		GetName:       getPolicyName,
		Type:          "policy",
		PrintHeader:   printHeaderFn,
	}

	commons.ExecuteOperation(ctx, client, &subjects, subjectOperator)
}
