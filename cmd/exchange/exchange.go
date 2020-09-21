package exchange

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/smartrecruiters/rabbitr/cmd/server"
	"github.com/urfave/cli"
)

func getExchanges(client *rabbithole.Client, vhost string) (*[]rabbithole.ExchangeInfo, error) {
	var exchanges []rabbithole.ExchangeInfo
	var err error
	if len(vhost) > 0 {
		exchanges, err = client.ListExchangesIn(vhost)
	} else {
		exchanges, err = client.ListExchanges()
	}
	return &exchanges, err
}

func getExchangeName(subject *interface{}) string {
	e := (*subject).(rabbithole.ExchangeInfo)
	return e.Name
}

func executeExchangeOperation(ctx *cli.Context, exchangeActionFn commons.SubjectActionFn, printHeaderFn commons.HeaderPrinterFn) {
	s := server.AskForServerSelection(ctx.String(commons.ServerName))
	vhost := ctx.String(commons.VHost)

	client := commons.GetRabbitClient(s)
	exchanges, err := getExchanges(client, vhost)
	commons.AbortIfError(err)

	subjects := commons.ConvertToSliceOfInterfaces(*exchanges)
	subjectOperator := commons.SubjectOperator{
		ExecuteAction: exchangeActionFn,
		GetName:       getExchangeName,
		Type:          "exchange",
		PrintHeader:   printHeaderFn,
	}
	commons.ExecuteOperation(ctx, client, &subjects, subjectOperator)
}
