package commons

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/urfave/cli"
)

const (
	// AllSubjects is a default filter for goevaluate that allows for returning all subjects
	// Useful for safe operations like list.
	AllSubjects = "1==1"
	// NoneOfTheSubjects is a default filter for goevaluate that allows for returning none subjects.
	// Useful for dangerous operations like pure or delete.
	NoneOfTheSubjects = "1!=1"
	// QueueFilterFields contains fields available when filtering queues
	QueueFilterFields = "queue.Name/Vhost/Durable/AutoDelete/Node/Status/Consumers/Policy/Messages/MessagesReady/Arguments (a map with string keys)"
	// ExchangeFilterFields contains fields available when filtering exchanges
	ExchangeFilterFields = "exchange.Name/Vhost/Type/Durable/AutoDelete/Internal/Arguments (a map with string keys)"
	// ShovelFilterFields contains fields available when filtering shovels
	ShovelFilterFields = "shovel.Name/Vhost"
)

// ServerFlag common server flag used in most of the commands
var ServerFlag cli.StringFlag

// DryRunFlag common dry run flag used in dangerous commands
var DryRunFlag cli.BoolFlag

// VHostFlag common virtual host flag used in most of the commands
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

// GetFilterFlag returns a flag describes a filter with given fields
func GetFilterFlag(defaultValue, availableFields string) cli.StringFlag {
	return cli.StringFlag{
		Name:  "filter, f",
		Value: defaultValue,
		Usage: fmt.Sprintf("Optional. Filter used to narrow list of subjects. It uses https://github.com/Knetic/govaluate engine. Fields available in filter: %s. Functions availble: getMapValueByKey(mapWithStringKeys, keyName)", availableFields),
	}
}

// AskIfValueEmpty prompts user for a value in case it was not provided beforehand
func AskIfValueEmpty(value, name string) string {
	value = strings.TrimSpace(value)
	if len(value) > 0 {
		return value
	}
	return strings.TrimSpace(AskWithValidator(value, name, NotEmptyValidator))
}

// AskForIntIf prompts user for an int value in case it was not provided beforehand
func AskForIntIf(testFn func(int) bool, value int, msg string) int {
	if !testFn(value) {
		return value
	}
	prompt := &survey.Input{
		Message: fmt.Sprintf("Please provide %s", msg),
	}
	customIntValidator := func(value interface{}) error {
		i, err := strconv.Atoi(value.(string))
		if err != nil || testFn(i) {
			return errors.New("please provide valid value")
		}
		return nil
	}

	err := survey.AskOne(prompt, &value, survey.WithValidator(customIntValidator))
	AbortIfError(err)
	return value
}

// AskForPasswordIfEmpty prompts user for a password in case it was not provided beforehand
func AskForPasswordIfEmpty(value, name string) string {
	value = strings.TrimSpace(value)
	if len(value) > 0 {
		return value
	}

	prompt := &survey.Password{
		Message: fmt.Sprintf("Please provide %s", name),
	}
	err := survey.AskOne(prompt, &value, survey.WithValidator(NotEmptyValidator))
	AbortIfError(err)
	return strings.TrimSpace(value)
}

// AskWithValidator prompts user for a value (using provided custom validator) in case it was not provided beforehand
func AskWithValidator(value, name string, validator survey.Validator) string {
	value = strings.TrimSpace(value)
	if validator(value) == nil {
		return value
	}
	prompt := &survey.Input{
		Message: fmt.Sprintf("Please provide %s", name),
	}
	err := survey.AskOne(prompt, &value, survey.WithValidator(validator))
	AbortIfError(err)
	return strings.TrimSpace(value)
}

// NotEmptyValidator returns validator that checks if provided value is not empty
func NotEmptyValidator(value interface{}) error {
	text := strings.TrimSpace(value.(string))
	if len(text) <= 0 {
		return errors.New("please provide not empty value")
	}
	return nil
}

// IsURL returns true if provided value is an url
func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Host != "" && (u.Scheme == "http" || u.Scheme == "https")
}

// IsURLValidator returns validator for checking urls
func IsURLValidator(value interface{}) error {
	text := strings.TrimSpace(value.(string))
	if !IsURL(text) {
		return errors.New("please provide valid URL")
	}
	return nil
}
