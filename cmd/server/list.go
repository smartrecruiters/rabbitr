package server

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli"

	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func showConfigurationCmd(ctx *cli.Context) error {
	cfg, _ := commons.GetApplicationConfig()
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	_, _ = fmt.Fprintln(w, "Server name\tAPI Url\tUser\tPassword")
	pass := "********"
	for name, s := range cfg.Servers {
		if ctx.Bool("show-passwords") {
			pass = s.Password
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", name, s.ApiURL, s.Username, pass)
	}
	_ = w.Flush()
	return nil
}
