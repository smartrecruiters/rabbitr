package shovel

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/smartrecruiters/rabbitr/cmd/server"
	"github.com/urfave/cli"
)

func getShovels(client *rabbithole.Client, vhost string) (*[]rabbithole.ShovelInfo, error) {
	var shovels []rabbithole.ShovelInfo
	var err error
	if len(vhost) > 0 {
		shovels, err = client.ListShovelsIn(vhost)
	} else {
		shovels, err = client.ListShovels()
	}
	return &shovels, err
}

func getShovelName(subject *interface{}) string {
	p := (*subject).(rabbithole.ShovelInfo)
	return p.Name
}

func executeShovelOperation(ctx *cli.Context, shovelActionFn commons.SubjectActionFn, printHeaderFn commons.HeaderPrinterFn) {
	s := server.AskForServerSelection(ctx.String(commons.ServerName))
	vhost := ctx.String(commons.VHost)

	client := commons.GetRabbitClient(s)
	shovels, err := getShovels(client, vhost)
	commons.AbortIfError(err)

	subjects := commons.ConvertToSliceOfInterfaces(*shovels)
	subjectOperator := commons.SubjectOperator{
		ExecuteAction: shovelActionFn,
		GetName:       getShovelName,
		Type:          "shovel",
		PrintHeader:   printHeaderFn,
	}

	commons.ExecuteOperation(ctx, client, &subjects, subjectOperator)
}
