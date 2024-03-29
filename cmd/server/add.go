package server

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"

	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func addServerCmd(ctx *cli.Context) {
	serverName := commons.AskIfValueEmpty(ctx.String(commons.ServerName), commons.ServerName)
	apiURL := commons.AskWithValidator(ctx.String("api-url"), "API url (for example: https://localhost:15672)", commons.IsURLValidator)
	amqpURL := commons.AskWithValidator(ctx.String("amqp-url"), "AMQP url (for example: amqp://localhost:5672)", commons.IsAmqpURLValidator)
	username := commons.AskIfValueEmpty(ctx.String("username"), "username")
	password := commons.AskForPasswordIfEmpty(ctx.String("password"), "password")

	fmt.Printf("Adding configuration for server %s:\n\t api url: %s\n\t amqp url: %s\n\t username: %s\n\t password: %s\n", serverName, apiURL, amqpURL, username, strings.Repeat("*", len(password)))

	cfg, err := commons.GetApplicationConfig(serverName)
	commons.AbortIfError(err)

	if cfg.Servers == nil {
		cfg.Servers = make(map[string]*commons.ServerCoordinates, 0)
	}

	cfg.Servers[serverName] = &commons.ServerCoordinates{
		APIURL:   apiURL,
		AmqpURL:  amqpURL,
		Username: username,
		Password: password,
	}

	err = commons.UpdateApplicationConfig(cfg, serverName)
	commons.AbortIfError(err)
}
