package commons

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	rabbithole "github.com/michaelklishin/rabbit-hole"
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

func PrintResponseBodyIfError(res *http.Response) {
	if res != nil && res.StatusCode >= 400 {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, res.Body)
		PrintIfError(err)
		fmt.Println(buf.String())
	}
}
