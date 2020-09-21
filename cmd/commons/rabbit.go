package commons

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func GetRabbitClient(serverName string) *rabbithole.Client {
	var exists bool
	var coordinates ServerCoordinates

	cfg, err := GetApplicationConfig()
	AbortIfError(err)

	if coordinates, exists = cfg.Servers[serverName]; !exists {
		fmt.Printf("configuration for server %s has not been found, please add it first via: `rabbitr server add` command", serverName)
		os.Exit(1)
	}
	client, err := rabbithole.NewClient(coordinates.ApiURL, coordinates.Username, coordinates.Password)
	AbortIfError(err)
	return client
}

func HandleGeneralResponse(messagePrefix string, res *http.Response) {
	if res == nil {
		return
	}
	Fprintf(os.Stdout, messagePrefix+", Response code: %d\t\n", res.StatusCode)
	PrintResponseBodyToWriterIfError(os.Stdout, res)
}

func HandleGeneralResponseWithWriter(w *tabwriter.Writer, res *http.Response) {
	if res == nil {
		return
	}
	Fprintf(w, "Response code: %d\t", res.StatusCode)
	PrintResponseBodyToWriterIfError(w, res)
}

func PrintResponseBodyToWriterIfError(w io.Writer, res *http.Response) {
	if res == nil || res.StatusCode < 400 {
		return
	}
	buf := new(strings.Builder)
	_, err := io.Copy(buf, res.Body)
	if err != nil {
		Fprintf(w, err.Error())
	}
	Fprintf(w, buf.String())
}
