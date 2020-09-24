package commons

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

var customFilterFunctions map[string]govaluate.ExpressionFunction

func init() {
	customFilterFunctions = map[string]govaluate.ExpressionFunction{
		"getMapValueByKey": func(args ...interface{}) (interface{}, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("invalid contains usage - expected: getMapValueByKey(mapWithStringKeys, 'keyName')")
			}
			mapWithStringKeys := args[0].(map[string]interface{})
			return mapWithStringKeys[args[1].(string)], nil
		},
	}
}

// GetCustomFilterFunctions returns custom goevaluate functions available for use during filtering
func GetCustomFilterFunctions() map[string]govaluate.ExpressionFunction {
	return customFilterFunctions
}
