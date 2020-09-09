package exchange

import (
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
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
	commons.AbortIfError(err)
	return &exchanges, err
}

func getExchangeName(subject *interface{}) string {
	e := (*subject).(rabbithole.ExchangeInfo)
	return e.Name
}

func executeExchangeOperation(ctx *cli.Context, exchangeActionFn commons.SubjectActionFn, printHeaderFn commons.HeaderPrinterFn) {
	server := ctx.String("server-name")
	vhost := ctx.String("vhost")

	client := commons.GetRabbitClient(server)
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
