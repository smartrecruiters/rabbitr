package server

import (
	"os"
	"text/tabwriter"

	"github.com/urfave/cli"

	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

func showConfigurationCmd(ctx *cli.Context) error {
	cfg, _ := commons.GetApplicationConfig("")
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	commons.Fprintln(w, "Server name\tApi Url\tAmqp Url\tUser\tPassword")
	pass := "********"

	for _, name := range cfg.GetServerNames() {
		s := cfg.Servers[name]
		if ctx.Bool("show-passwords") {
			pass = s.Password
		}
		commons.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", name, s.APIURL, s.AmqpURL, s.Username, pass)
	}
	_ = w.Flush()
	return nil
}
