package connection

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/smartrecruiters/rabbitr/cmd/server"
	"github.com/urfave/cli"
)

type ConnectionInfo struct {
	ID    string
	Name  string
	Vhost string
}

func getConnections(client *rabbithole.Client, vhost string) (*[]ConnectionInfo, error) {
	connections, err := client.ListConnections()
	connectionInfos := make([]ConnectionInfo, 0)
	for _, connection := range connections {
		clientProvidedName := connection.ClientProperties["connection_name"]
		if clientProvidedName == nil {
			clientProvidedName = "not-defined"
		}
		if len(vhost) <= 0 || vhost == connection.Vhost {
			connectionInfos = append(connectionInfos, ConnectionInfo{ID: connection.Name, Name: clientProvidedName.(string), Vhost: connection.Vhost})
		}
	}
	return &connectionInfos, err
}

func getConnectionName(subject *interface{}) string {
	c := (*subject).(ConnectionInfo)
	return c.Name
}

func executeConnectionOperation(ctx *cli.Context, connectionActionFn commons.SubjectActionFn, headerPrinterFn commons.HeaderPrinterFn) {
	s := server.AskForServerSelection(ctx.String(commons.ServerName))
	vhost := ctx.String(commons.VHost)

	client := commons.GetRabbitClient(s)
	queues, err := getConnections(client, vhost)
	commons.AbortIfError(err)

	subjects := commons.ConvertToSliceOfInterfaces(*queues)
	subjectOperator := commons.SubjectOperator{
		ExecuteAction: connectionActionFn,
		GetName:       getConnectionName,
		Type:          "connection",
		PrintHeader:   headerPrinterFn,
	}
	commons.ExecuteOperation(ctx, client, &subjects, subjectOperator)
}
