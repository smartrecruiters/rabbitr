package server

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func AskForServerSelection(server string) string {
	server = strings.TrimSpace(server)
	cfg := commons.GetCachedApplicationConfig()
	serverNames := cfg.GetServerNames()
	if len(server) <= 0 || !commons.Contains(serverNames, server) {
		prompt := &survey.Select{
			Message:  "Please choose a server that you wish to act upon:",
			Options:  serverNames,
			PageSize: 20,
		}
		err := survey.AskOne(prompt, &server, nil)
		commons.AbortIfError(err)
	}
	return strings.TrimSpace(server)
}
