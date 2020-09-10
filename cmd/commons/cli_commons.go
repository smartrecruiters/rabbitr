package commons

import (
	"errors"
	"fmt"
	"gopkg.in/AlecAivazis/survey.v2"
	"net/url"
	"strings"

	"github.com/urfave/cli"
)

const AllSubjects = "1==1"
const NoneOfTheSubjects = "1!=1"
const QueueFilterFields = "queue.Name/Vhost/Durable/AutoDelete/Node/Status/Consumers/Policy/Messages/MessagesReady/Arguments (a map with string keys)"
const ExchangeFilterFields = "exchange.Name/Vhost/Type/Durable/AutoDelete/Internal/Arguments (a map with string keys)"

var ServerFlag cli.StringFlag
var DryRunFlag cli.BoolFlag
var VHostFlag cli.StringFlag

func init() {
	ServerFlag = cli.StringFlag{
		Name:  "server, server-name, s",
		Value: "",
		Usage: "Server name used to perform given operation",
	}
	DryRunFlag = cli.BoolFlag{
		Name:  "dry-run",
		Usage: "Dry run, will only print subjects that would be acted upon but without actually performing action on them",
	}
	VHostFlag = cli.StringFlag{
		Name:  "vhost, v",
		Usage: "Optional. Virtual host used to narrow list of subjects. If not provided all vhosts are considered",
	}
}

func GetFilterFlag(defaultValue, availableFields string) cli.StringFlag {
	return cli.StringFlag{
		Name:  "filter, f",
		Value: defaultValue,
		Usage: fmt.Sprintf("Optional. Filter used to narrow list of subjects. It uses https://github.com/Knetic/govaluate engine. Fields available in filter: %s. Functions availble: getMapValueByKey(mapWithStringKeys, keyName)", availableFields),
	}
}

func AskIfValueEmpty(value, name string) string {
	value = strings.TrimSpace(value)
	if len(value) > 0 {
		return value
	}
	return strings.TrimSpace(AskWithValidator(value, name, NotEmptyValidator))
}

func AskWithValidator(value, name string, validator survey.Validator) string {
	value = strings.TrimSpace(value)
	prompt := &survey.Input{
		Message: fmt.Sprintf("Please provide %s", name),
	}
	err := survey.AskOne(prompt, &value, survey.WithValidator(validator))
	AbortIfError(err)
	return strings.TrimSpace(value)
}

func NotEmptyValidator(value interface{}) error {
	text := strings.TrimSpace(value.(string))
	if len(text) <= 0 {
		return errors.New("please provide not empty value")
	}
	return nil
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func IsUrlValidator(value interface{}) error {
	text := strings.TrimSpace(value.(string))
	if !IsUrl(text) {
		return errors.New("please provide valid URL")
	}
	return nil
}
