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

// GetRabbitClient returns rabbit client initialized with a provider server coordinates
func GetRabbitClient(serverName string) *rabbithole.Client {
	var exists bool
	var coordinates ServerCoordinates

	cfg, err := GetApplicationConfig()
	AbortIfError(err)

	if coordinates, exists = cfg.Servers[serverName]; !exists {
		fmt.Printf("configuration for server %s has not been found, please add it first via: `rabbitr server add` command", serverName)
		os.Exit(1)
	}
	client, err := rabbithole.NewClient(coordinates.APIURL, coordinates.Username, coordinates.Password)
	AbortIfError(err)
	return client
}

// HandleGeneralResponse handles general response from rabbit client.
// Prints response code and error details if able.
func HandleGeneralResponse(messagePrefix string, res *http.Response) {
	if res == nil {
		return
	}
	Fprintf(os.Stdout, messagePrefix+", Response code: %d\t\n", res.StatusCode)
	PrintResponseBodyToWriterIfError(os.Stdout, res)
}

// HandleGeneralResponseWithWriter handles general response from rabbit client.
// Prints response code and error details if able.
// Writes output to the provided writer.
func HandleGeneralResponseWithWriter(w *tabwriter.Writer, res *http.Response) {
	if res == nil {
		return
	}
	Fprintf(w, "Response code: %d\t", res.StatusCode)
	PrintResponseBodyToWriterIfError(w, res)
}

// PrintResponseBodyToWriterIfError prints response body to provided writer
// for responses with code >= 400 to obtain as much error details as possible.
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
