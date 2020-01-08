package commons

import (
	"log"
	"runtime"
	"strings"
)

func AbortIfTrue(condition bool, message string) {
	if condition {
		log.Fatal(message)
	}
}

func AbortIfError(err error, messages ...interface{}) {
	if err != nil {
		if len(messages) > 0 {
			log.Fatal(messages...)
		} else {
			log.Fatal(err)
		}
	}
}

func AbortIfErrorWithMsg(msg string, err error) {
	if err != nil {
		log.Fatalf(msg, err.Error())
	}
}

func PrintIfTrue(condition bool, msg string) {
	if condition {
		log.Println(msg)
	}
}

func PrintIfErrorWithMsg(msg string, err error) {
	if err != nil {
		log.Printf(msg+"\n", err.Error())
	}
}

func PrintIfError(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}

func IsOSX() bool {
	return strings.Contains(runtime.GOOS, "darwin")
}
