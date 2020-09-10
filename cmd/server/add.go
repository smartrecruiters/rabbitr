package server

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func addServerCmd(ctx *cli.Context) {
	serverName := commons.AskIfValueEmpty(ctx.String("server-name"), "server-name")
	apiUrl := commons.AskWithValidator(ctx.String("url"), "API url (for example: https://localhost:15672)", commons.IsUrlValidator)
	username := commons.AskIfValueEmpty(ctx.String("username"), "username")
	password := commons.AskIfValueEmpty(ctx.String("password"), "password")

	fmt.Printf("Adding configuration for server %s:\n\t api url: %s\n\t username: %s\n\t password: %s\n", serverName, apiUrl, username, password)

	cfg, err := commons.GetApplicationConfig()
	commons.AbortIfError(err)

	if cfg.Servers == nil {
		cfg.Servers = make(map[string]commons.ServerCoordinates, 0)
	}

	cfg.Servers[serverName] = commons.ServerCoordinates{
		ApiURL:   apiUrl,
		Username: username,
		Password: password,
	}

	err = commons.UpdateApplicationConfig(cfg)
	commons.AbortIfError(err)
}
