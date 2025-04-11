package connection

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
	"github.com/smartrecruiters/rabbitr/cmd/rabbit"
	"github.com/smartrecruiters/rabbitr/cmd/server"
	"github.com/urfave/cli"
)

// ConnInfo contains details about a rabbitmq connection
type ConnInfo struct {
	ID               string
	Name             string
	User             string
	Vhost            string
	ClientProperties map[string]interface{}
}

func getConnections(client *rabbithole.Client, vhost string) (*[]ConnInfo, error) {
	connections, err := client.ListConnections()
	connectionInfos := make([]ConnInfo, 0)
	for _, connection := range connections {
		clientProvidedName := connection.ClientProperties["connection_name"]
		if clientProvidedName == nil {
			clientProvidedName = "not-defined"
		}
		if len(vhost) <= 0 || vhost == connection.Vhost {
			connectionInfos = append(connectionInfos, ConnInfo{ID: connection.Name,
				Name:             clientProvidedName.(string),
				Vhost:            connection.Vhost,
				User:             connection.User,
				ClientProperties: connection.ClientProperties})
		}
	}
	return &connectionInfos, err
}

func getConnectionName(subject *interface{}) string {
	c := (*subject).(ConnInfo)
	return c.Name
}

func executeConnectionOperation(ctx *cli.Context, connectionActionFn commons.SubjectActionFn, headerPrinterFn commons.HeaderPrinterFn) {
	s := server.AskForServerSelection(ctx.String(commons.ServerName))
	vhost := ctx.String(commons.VHost)

	client := rabbit.GetRabbitClient(s)
	connections, err := getConnections(client, vhost)
	commons.AbortIfError(err)

	subjects := commons.ConvertToSliceOfInterfaces(*connections)
	subjectOperator := commons.SubjectOperator{
		ExecuteAction: connectionActionFn,
		GetName:       getConnectionName,
		Type:          "connection",
		PrintHeader:   headerPrinterFn,
	}
	commons.ExecuteOperation(ctx, client, &subjects, subjectOperator)
}
