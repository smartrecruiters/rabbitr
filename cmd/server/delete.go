package server

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func deleteServerConfigCmd(ctx *cli.Context) {
	serverCfgToRemove := AskForServerSelection(ctx.String(commons.ServerName))

	cfg, err := commons.GetApplicationConfig()
	commons.AbortIfError(err)

	fmt.Printf("Removing configuration for %s server\n", serverCfgToRemove)
	delete(cfg.Servers, serverCfgToRemove)

	err = commons.UpdateApplicationConfig(cfg)
	commons.AbortIfError(err)
}
